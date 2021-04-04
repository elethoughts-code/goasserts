package fsbuilder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Common interface {
	Parent() D
	Root() D
	Remove()
	Name() string
}

type D interface {
	Common
	Dir(name string, perm os.FileMode) D
	File(name string, flag int, perm os.FileMode) F
	RemoveAll()
}

type F interface {
	Common
	Write(content []byte) F
	WriteString(content string) F
	WriteStringDedent(content string) F
}

func TmpDir(dir, pattern string) D {
	name, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		panic(err)
	}
	return &dirBuilder{name: name, parent: nil}
}

type dirBuilder struct {
	name   string
	parent D
}

func (d *dirBuilder) Name() string {
	return d.name
}

func (d *dirBuilder) Root() D {
	if d.parent == nil {
		return d
	}
	return d.parent.Root()
}

func (d *dirBuilder) Remove() {
	if err := os.Remove(d.name); err != nil {
		panic(err)
	}
}

func (d *dirBuilder) RemoveAll() {
	if err := os.RemoveAll(d.name); err != nil {
		panic(err)
	}
}

func (d *dirBuilder) Dir(name string, perm os.FileMode) D {
	name = filepath.Join(d.Name(), name)
	if err := os.Mkdir(name, perm); err != nil {
		panic(err)
	}
	return &dirBuilder{name: name, parent: d}
}

func (d *dirBuilder) File(name string, flag int, perm os.FileMode) F {
	name = filepath.Join(d.Name(), name)
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()
	return &fileBuilder{name: name, parent: d}
}

func (d *dirBuilder) Parent() D {
	return d.parent
}

type fileBuilder struct {
	name   string
	parent D
}

func (f *fileBuilder) Name() string {
	return f.name
}

func (f *fileBuilder) Root() D {
	return f.Parent().Root()
}

func (f *fileBuilder) Remove() {
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
}

func (f *fileBuilder) Parent() D {
	return f.parent
}

func (f *fileBuilder) Write(content []byte) F {
	if err := ioutil.WriteFile(f.name, content, 0600); err != nil {
		panic(err)
	}
	return f
}

func (f *fileBuilder) WriteString(content string) F {
	return f.Write([]byte(content))
}

var prefixSpacesRegexp = regexp.MustCompile(`(?m)(^[ \t]*)(?:[^ \t\n])`)

func (f *fileBuilder) WriteStringDedent(content string) F {
	prefixSpaces := ""

	prefixes := prefixSpacesRegexp.FindAllStringSubmatch(content, -1)

l:
	for i, p := range prefixes {
		switch {
		case i == 0:
			prefixSpaces = p[1]
		case strings.HasPrefix(p[1], prefixSpaces):
			// current is prefixed by selected prefix
			continue
		case strings.HasPrefix(prefixSpaces, p[1]):
			// current is a shorter common prefix
			prefixSpaces = p[1]
		default:
			// no common prefix
			prefixSpaces = ""
			break l
		}
	}

	if prefixSpaces != "" {
		content = regexp.MustCompile(fmt.Sprintf("(?m)^%s", prefixSpaces)).ReplaceAllString(content, "")
	}

	return f.WriteString(content)
}
