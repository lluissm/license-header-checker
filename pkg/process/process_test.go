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
)

// fileInfoStub implements os.FileInfo
type dirEntryStub struct {
	isDir bool
}

func (d *dirEntryStub) Name() string               { return "name" }
func (d *dirEntryStub) IsDir() bool                { return d.isDir }
func (d *dirEntryStub) Type() os.FileMode          { return 0 }
func (d *dirEntryStub) Info() (os.FileInfo, error) { return nil, nil }

type ioHandlerStub struct {
	pathsToWalk               []string
	isDir                     bool
	errorReplacingFileContent bool
	errorWalkingPath          bool
}

func (s *ioHandlerStub) ReplaceFileContent(filePath string, fileContent string) error {
	if s.errorReplacingFileContent {
		return errors.New("error")
	}
	return nil
}

func (s *ioHandlerStub) ReadFile(filename string) ([]byte, error) {
	switch filename {
	case "license.txt":
		return []byte(testTargetLicenseHeader), nil
	case "file_no_license.cpp":
		return []byte(testFileWithoutLicense), nil
	case "file_good_license.cpp":
		return []byte(testFileWithTargetLicense), nil
	case "file_old_license.cpp":
		return []byte(testFileWithDifferentLicense), nil
	default:
		return nil, errors.New("file does not exist")
	}
}

func (s *ioHandlerStub) WalkDir(path string, walkDirFn fs.WalkDirFunc) error {
	fileInfo := new(dirEntryStub)
	fileInfo.isDir = s.isDir
	for _, path := range s.pathsToWalk {
		if s.errorWalkingPath {
			if err := walkDirFn(path, fileInfo, errors.New("error")); err != nil {
				return err
			}
		} else {
			if err := walkDirFn(path, fileInfo, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestFileLicenseOk(t *testing.T) {
	filePath := "main.cpp"
	license := testTargetLicenseHeader
	fileContent := testFileWithTargetLicense
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseOk)
}

func TestFileAddLicense(t *testing.T) {
	filePath := "main.cpp"
	license := testTargetLicenseHeader
	fileContent := testFileWithoutLicense
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	options.Add = false
	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == SkippedAdd)

	options.Add = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseAdded)

	options.Add = true
	handler.errorReplacingFileContent = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == OperationError)
}

func TestFileReplaceLicense(t *testing.T) {
	filePath := "main.cpp"
	license := testTargetLicenseHeader
	fileContent := testFileWithDifferentLicense
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	options.Replace = false
	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == SkippedReplace)

	options.Replace = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseReplaced)

	options.Replace = true
	handler.errorReplacingFileContent = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == OperationError)
}

func TestFilesSuccess(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	handler.pathsToWalk = []string{"file_no_license.cpp", // license to add
		"file_good_license.cpp", // license is ok
		"file_old_license.cpp",  // license to replace
		"file_old_license.h",    // extension to ignore
		"ignore/file.cpp"}       // path to ignore

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[LicenseOk]) == 1)
	assert.True(t, len(stats.Files[LicenseReplaced]) == 1)
	assert.True(t, len(stats.Files[LicenseAdded]) == 1)
}

func TestFilesErrorReadingLicense(t *testing.T) {
	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "wrong_license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	handler.pathsToWalk = []string{"file_no_license.cpp"}

	_, err := Files(options, handler)
	assert.NotNil(t, err)
}

func TestFilesDoesNotCountDir(t *testing.T) {
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

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[LicenseAdded]) == 0)
}

func TestFilesErrorReadingFile(t *testing.T) {

	handler := new(ioHandlerStub)
	options := &Options{
		Add:         true,
		Replace:     true,
		Path:        "source",
		LicensePath: "license.txt",
		Extensions:  []string{".cpp"},
		IgnorePaths: []string{"ignore"},
	}

	handler.pathsToWalk = []string{"file_does_not_exist.cpp"}

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[OperationError]) == 1)
}

func TestFilesErrorSentByWalk(t *testing.T) {
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
	handler.errorWalkingPath = true

	stats, err := Files(options, handler)
	assert.Nil(t, err)
	assert.True(t, len(stats.Files[OperationError]) == 1)
}
