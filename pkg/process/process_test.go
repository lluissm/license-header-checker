/* MIT License

Copyright (c) 2022 Lluis Sanchez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package process

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// dirEntryMock mocks fs.DirEntry
type dirEntryMock struct {
	mock.Mock
}

func (d *dirEntryMock) IsDir() bool {
	args := d.Called()
	return args.Get(0).(bool)
}
func (d *dirEntryMock) Name() string               { return "" }
func (d *dirEntryMock) Type() os.FileMode          { return 0 }
func (d *dirEntryMock) Info() (os.FileInfo, error) { return nil, nil }

// ioHandlerStub mocks ioHandle
type ioHandlerStub struct {
	mock.Mock
	pathsToWalk      []string
	isDir            bool
	errorWalkingPath bool
}

func (s *ioHandlerStub) ReplaceFileContent(filePath string, fileContent string) error {
	args := s.Called(filePath, fileContent)
	return args.Error(0)
}

func (s *ioHandlerStub) ReadFile(filename string) ([]byte, error) {
	args := s.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *ioHandlerStub) WalkDir(path string, walkDirFn fs.WalkDirFunc) error {
	args := s.Called(path, walkDirFn)

	dirEntry := &dirEntryMock{}
	dirEntry.On("IsDir").Return(s.isDir)
	for _, path := range s.pathsToWalk {
		var errSent error
		if s.errorWalkingPath {
			errSent = errors.New("error")
		}
		if err := walkDirFn(path, dirEntry, errSent); err != nil {
			return nil
		}
	}

	return args.Error(0)
}

func TestFile_LicenseOk(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	op := File("main.cpp", testFileWithTargetLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == LicenseOk)
}

func TestFile_AddLicense(t *testing.T) {
	fileName := "main.cpp"
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	// Return SkippedAdd if the file does not contain a license BUT options.Add is false
	options.Add = false
	op := File(fileName, testFileWithoutLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == SkippedAdd)

	// Return LicenseAdded if the file does not contain a license and options.Add is true
	options.Add = true
	handler.On("ReplaceFileContent", fileName, testFileWithTargetLicense).Return(nil).Once()
	op = File(fileName, testFileWithoutLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == LicenseAdded)
	handler.AssertExpectations(t)

	// Return OperationError if there was an error saving the file
	options.Add = true
	handler.On("ReplaceFileContent", fileName, testFileWithTargetLicense).Return(errors.New("error")).Once()
	op = File(fileName, testFileWithoutLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == OperationError)
}

func TestFile_ReplaceLicense(t *testing.T) {
	fileName := "main.cpp"
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	// Return SkippedReplace if the file does not contain a license BUT options.Add is false
	options.Replace = false
	op := File(fileName, testFileWithDifferentLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == SkippedReplace)

	// Return LicenseReplaced if the file does not contain a license and options.Add is true
	options.Replace = true
	handler.On("ReplaceFileContent", fileName, testFileWithTargetLicense).Return(nil).Once()
	op = File(fileName, testFileWithDifferentLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == LicenseReplaced)
	handler.AssertExpectations(t)

	// Return OperationError if there was an error saving the file
	options.Replace = true
	handler.On("ReplaceFileContent", fileName, testFileWithTargetLicense).Return(errors.New("error")).Once()
	op = File(fileName, testFileWithDifferentLicense, testTargetLicenseHeader, options, handler)
	assert.True(t, op == OperationError)
}

func TestFiles_Success(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	// Prepare mock ioHandler to walk through all the files
	handler.pathsToWalk = []string{
		"file_no_license.cpp",   // license to add
		"file_good_license.cpp", // license is ok
		"file_old_license.cpp",  // license to replace
		"file_old_license.h",    // extension to ignore
		"ignore/file.cpp"}       // path to ignore
	handler.On("WalkDir", options.Path, mock.Anything).Return(nil).Once()

	// ReadFile should return the files that should be read (not ignored)
	handler.On("ReadFile", "license.txt").Return([]byte(testTargetLicenseHeader), nil).Once()
	handler.On("ReadFile", "file_no_license.cpp").Return([]byte(testFileWithoutLicense), nil).Once()
	handler.On("ReadFile", "file_good_license.cpp").Return([]byte(testFileWithTargetLicense), nil).Once()
	handler.On("ReadFile", "file_old_license.cpp").Return([]byte(testFileWithDifferentLicense), nil).Once()

	// ReplaceFileContent should be called for the 2 files where license will be added/replaced
	handler.On("ReplaceFileContent", "file_no_license.cpp", mock.Anything).Return(nil).Times(1)
	handler.On("ReplaceFileContent", "file_old_license.cpp", mock.Anything).Return(nil).Times(1)

	// Execute and validate stats
	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[LicenseOk]) == 1)
	assert.True(t, len(stats.Files[LicenseReplaced]) == 1)
	assert.True(t, len(stats.Files[LicenseAdded]) == 1)

	handler.AssertExpectations(t)
}

func TestFiles_ErrorReadingLicense(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	// Return error while reading the license as if it did not exist
	handler.On("ReadFile", "license.txt").Return([]byte{}, errors.New("error")).Once()

	// Assert that Files will return an error
	_, err := Files(options, handler)
	assert.NotNil(t, err)

	handler.AssertExpectations(t)
}

func TestFiles_ErrorReadingFile(t *testing.T) {

	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	handler.pathsToWalk = []string{"file.cpp"}

	// Prepare mock ioHandler to walk through all the files
	handler.pathsToWalk = []string{"file.cpp"}
	handler.On("WalkDir", options.Path, mock.Anything).Return(nil).Once()

	// ReadFile should return the license
	handler.On("ReadFile", "license.txt").Return([]byte(testTargetLicenseHeader), nil).Once()

	// ReadFile should return error while reading the file
	handler.On("ReadFile", "file.cpp").Return([]byte{}, errors.New("error")).Once()

	// Assert that Files will return nil and operation will be OperationError
	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[OperationError]) == 1)

	handler.AssertExpectations(t)
}

func TestFiles_ErrorSentByWalk(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	// Prepare mock ioHandler to return an error on WalkDir
	handler.On("WalkDir", options.Path, mock.Anything).Return(nil).Once()
	handler.pathsToWalk = []string{"some_file.cpp"}
	handler.errorWalkingPath = true

	// ReadFile should return the license
	handler.On("ReadFile", "license.txt").Return([]byte(testTargetLicenseHeader), nil).Once()

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[OperationError]) == 1)

	handler.AssertExpectations(t)
}

func TestFiles_DoesNotCountDir(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	handler.pathsToWalk = []string{"file_no_license.cpp"}
	handler.isDir = true

	// ReadFile should return the license
	handler.On("ReadFile", "license.txt").Return([]byte(testTargetLicenseHeader), nil).Once()

	// Prepare mock ioHandler to return an error on WalkDir
	handler.On("WalkDir", options.Path, mock.Anything).Return(nil).Once()

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[LicenseAdded]) == 0)

	handler.AssertExpectations(t)
}
