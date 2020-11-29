package fsbuilder_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	fsBuilder "github.com/elethoughts-code/goasserts/fs_builder"
)

func Test_should_create_multiple_folder_hierarchy(t *testing.T) {
	assert := assertion.New(t)

	tmp := fsBuilder.TmpDir("", ".")
	tmp.Dir("d1", 0755).Dir("d2", 0755).Dir("d31", 0755).Parent().Dir("d32", 0755)

	assert.That(filepath.Join(tmp.Name(), "d1")).FileExists()
	assert.That(filepath.Join(tmp.Name(), "d1", "d2")).FileExists()
	assert.That(filepath.Join(tmp.Name(), "d1", "d2", "d31")).FileExists()
	assert.That(filepath.Join(tmp.Name(), "d1", "d2", "d32")).FileExists()
}

func Test_Builder_should_make_tmp_folders_and_files(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")
	subFolder1 := tmpFolder.Dir("my_sub_folder", 0755)

	file1 := tmpFolder.File("my_file_1", os.O_CREATE, 0755)
	file2 := subFolder1.File("my_file_2", os.O_CREATE, 0755)

	file2.Write([]byte("Hello world !"))

	// Then
	assert.That(tmpFolder.Name()).HasPrefix(filepath.Join(os.TempDir(), "my_folder_1"))
	assert.That(subFolder1.Name()).IsEq(filepath.Join(tmpFolder.Name(), "my_sub_folder"))
	assert.That(file1.Name()).IsEq(filepath.Join(tmpFolder.Name(), "my_file_1"))
	assert.That(file2.Name()).IsEq(filepath.Join(subFolder1.Name(), "my_file_2"))

	f2Content, err := ioutil.ReadFile(file2.Name())
	assert.That(err).IsNil()
	assert.That(string(f2Content)).IsEq("Hello world !")

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_Builder_should_write_file_dedent(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")
	file1 := tmpFolder.
		File("my_file_1", os.O_CREATE, 0755).
		WriteStringDedent(`
						a:
						  b: some text
						  c: [1,2,3]
						  d:
						    - 1
						    - 2`)
	// Then
	content, err := ioutil.ReadFile(file1.Name())
	assert.That(err).IsNil()
	assert.That(string(content)).IsEq(`
a:
  b: some text
  c: [1,2,3]
  d:
    - 1
    - 2`)
	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_Builder_should_write_file_dedent_2(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")
	file1 := tmpFolder.
		File("my_file_1", os.O_CREATE, 0755).
		WriteStringDedent(`
						a:
						  b: some text
					c: [1,2,3]
						  d:
						    - 1
						    - 2`)
	// Then
	content, err := ioutil.ReadFile(file1.Name())
	assert.That(err).IsNil()
	assert.That(string(content)).IsEq(`
	a:
	  b: some text
c: [1,2,3]
	  d:
	    - 1
	    - 2`)
	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_Builder_should_not_dedent_when_no_consistent_leading_spaces(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")
	file1 := tmpFolder.
		File("my_file_1", os.O_CREATE, 0755).
		WriteStringDedent(`
							a:
						  b: some text
                          c: [1,2,3]
						  d: # there are spaces and not tabs
						    - 1
						    - 2`)
	// Then
	content, err := ioutil.ReadFile(file1.Name())
	assert.That(err).IsNil()
	assert.That(string(content)).IsEq(`
							a:
						  b: some text
                          c: [1,2,3]
						  d: # there are spaces and not tabs
						    - 1
						    - 2`)
	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_TmpDir_should_panic_when_trying_to_create_folder_into_non_existing_one(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()
	fsBuilder.TmpDir("/not_here", "my_folder_1")
}

func Test_Remove_should_panic_when_already_removed_folder(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()
	tmp := fsBuilder.TmpDir("", "my_folder_1")
	tmp.Remove()
	tmp.Remove()
}

func Test_Remove_should_panic_when_already_removed_file(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()
	tmp := fsBuilder.TmpDir("", "my_folder_1")
	file := tmp.File("my_file", os.O_CREATE, 0755)
	file.Remove()
	file.Remove()

	t.Cleanup(tmp.RemoveAll)
}

func Test_Dir_should_panic_when_already_created(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrExist)
	}()
	tmp := fsBuilder.TmpDir("", "folder")
	tmp.File("my_dir", os.O_CREATE, 0755).Parent().Dir("my_dir", 0755)
	t.Cleanup(tmp.RemoveAll)
}

func Test_File_should_panic_when_not_creating_and_non_existing_file(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()
	tmp := fsBuilder.TmpDir("", "folder")
	tmp.File("my_file", os.O_RDONLY, 0755)
	t.Cleanup(tmp.RemoveAll)
}

func Test_File_should_panic_when_writing_to_removed_file(t *testing.T) {
	assert := assertion.New(t)
	defer func() {
		r := recover()
		assert.That(r).IsError(os.ErrNotExist)
	}()
	tmp := fsBuilder.TmpDir("", "folder")
	f := tmp.File("my_file", os.O_CREATE, 0755)
	tmp.RemoveAll()
	f.WriteString("abc")
}
