package assertion_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/elethoughts-code/goasserts/assertion"
	fsBuilder "github.com/elethoughts-code/goasserts/fs_builder"
	mocks "github.com/elethoughts-code/goasserts/mocks/assertion"
	"github.com/golang/mock/gomock"
)

func Test_File_exists_should_check_file_existence(t *testing.T) {
	assert := assertion.New(t)

	// When
	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	file4 := tmpFolder.Dir("my_sub_folder", 0755).
		File("file_1", os.O_CREATE, 0755).Parent().
		File("file_2", os.O_CREATE, 0755).Parent().
		File("file_3", os.O_CREATE, 0755).Root().
		File("file_4", os.O_CREATE, 0755)

	// Then

	assert.That(filepath.Join(tmpFolder.Name(), "my_sub_folder")).FileExists()
	assert.That(filepath.Join(tmpFolder.Name(), "my_sub_folder", "file_1")).FileExists()
	assert.That(filepath.Join(tmpFolder.Name(), "my_sub_folder", "file_2")).FileExists()
	assert.That(filepath.Join(tmpFolder.Name(), "my_sub_folder", "file_3")).FileExists()
	assert.That(filepath.Join(tmpFolder.Name(), "file_4")).FileExists()

	// When
	file4.Remove()

	// Then
	assert.That(filepath.Join(tmpFolder.Name(), "file_4")).Not().FileExists()

	// Clean up
	t.Cleanup(tmpFolder.Root().RemoveAll)
}

func Test_Fs_Matchers_should_fail(t *testing.T) {
	// Mock preparation
	ctrl := gomock.NewController(t)

	tmpFolder := fsBuilder.TmpDir("", "my_folder_1")

	testEntries := []struct {
		assertFunc func(assert assertion.Assert)
		errLog     string
	}{
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That("/my_file").FileExists()
			},
			errLog: "\nFile /my_file do not exists",
		},
		{
			assertFunc: func(assert assertion.Assert) {
				assert.That(tmpFolder.Name()).Not().FileExists()
			},
			errLog: fmt.Sprintf("\nFile %s should not exists", tmpFolder.Name()),
		},
	}

	for _, entry := range testEntries {
		// Given
		tMock := mocks.NewMockPublicTB(ctrl)
		assert := assertion.New(tMock)

		// Expectation
		tMock.EXPECT().Helper().AnyTimes()
		tMock.EXPECT().Error(entry.errLog)

		// When
		entry.assertFunc(assert)
	}
}
