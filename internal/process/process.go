/* MIT License

Copyright (c) 2020 Lluis Sanchez

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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/lsm-dev/license-header-checker/internal/header"
)

type (
	// Action performed when processing a file
	Action int

	// Operation is the result of processing one file
	Operation struct {
		Action Action
		Path   string
	}

	// Stats is the result of processing multiple files
	Stats struct {
		Operations []*Operation
		ElapsedMs  int64
	}

	// Options to be followed during processing
	Options struct {
		Add         bool
		Replace     bool
		Path        string
		LicensePath string
		Extensions  []string
		IgnorePaths []string
	}
)

const (
	// SkippedAdd means that the file had no license but the new one was not added. Missing -a flasg
	SkippedAdd Action = iota
	// SkippedReplace means that the file had a different license but it was not replaced with the target one. Missing -r flag
	SkippedReplace
	// LicenseOk means that the license was OK
	LicenseOk
	// LicenseAdded means that the target license was added to the file
	LicenseAdded
	// LicenseReplaced means that the license was replaced by the target one
	LicenseReplaced
	// OperationError means there was an error with one of the files
	OperationError
)

// Files processes all files in the path that match the options
func Files(options *Options) (*Stats, error) {

	license, err := getLicense(options.LicensePath)
	if err != nil {
		return nil, err
	}

	channel := make(chan *Operation)
	start := time.Now()
	files := 0

	err = filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			files++
			channel <- &Operation{
				Action: OperationError,
				Path:   path,
			}
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if shouldIgnorePath(path, options.IgnorePaths) {
			return nil
		}
		if shouldIgnoreExtension(path, options.Extensions) {
			return nil
		}

		files++
		go func() {
			action := File(path, license, options)
			channel <- &Operation{
				Action: action,
				Path:   path,
			}
		}()

		return nil
	})

	operations := []*Operation{}
	for i := 0; i < files; i++ {
		operations = append(operations, <-channel)
	}

	elapsedTime := time.Since(start)
	return &Stats{
		Operations: operations,
		ElapsedMs:  elapsedTime.Milliseconds(),
	}, err
}

// File processes one file
func File(filePath string, license string, options *Options) Action {

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return OperationError
	}

	fileContent := string(data)

	if strings.Contains(fileContent, strings.TrimSpace(license)) {
		return LicenseOk
	}

	if header.ContainsLicense(fileContent) {
		if options.Replace {
			if err := replaceLicense(filePath, fileContent, license); err != nil {
				return OperationError
			}
			return LicenseReplaced
		}
		return SkippedReplace
	}

	if options.Add {
		if err := addLicense(filePath, fileContent, license); err != nil {
			return OperationError
		}
		return LicenseAdded
	}
	return SkippedAdd
}

// getLicense reads the license from the path
func getLicense(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// shouldIgnore returns true if the path matches any of the paths to ignore
func shouldIgnorePath(path string, ignorePaths []string) bool {
	pathSegments := strings.Split(path, string(os.PathSeparator))
	for _, ignorePath := range ignorePaths {
		ignorePathSegments := strings.Split(ignorePath, string(os.PathSeparator))
		size := len(ignorePathSegments)
		lastSegment := len(pathSegments) - size
		for i := 0; i <= lastSegment; i++ {
			if reflect.DeepEqual(pathSegments[i:i+size], ignorePathSegments) {
				return true
			}
		}
	}
	return false
}

// shouldIgnoreExtension returns false only if the file's extension is one of the provided ones
func shouldIgnoreExtension(path string, extensions []string) bool {
	fileExtension := filepath.Ext(path)
	for _, ext := range extensions {
		if fileExtension == ext {
			return false
		}
	}
	return true
}

func addLicense(filePath string, fileContent string, license string) error {
	newFileContent := header.Insert(fileContent, license)
	return replaceFile(filePath, newFileContent)
}

func replaceLicense(filePath string, fileContent string, license string) error {
	newFileContent := header.Replace(fileContent, license)
	return replaceFile(filePath, newFileContent)
}

// replaceFile removes the file and create a new one with the specified content
func replaceFile(filePath string, content string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed deleting the file: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	return nil
}
