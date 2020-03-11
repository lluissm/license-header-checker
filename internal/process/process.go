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
	"os"
	"path/filepath"
	"strings"
	"time"
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

type ioHandle interface {
	ReplaceFileContent(filePath string, content string) error
	Walk(string, filepath.WalkFunc) error
	ReadFile(string) ([]byte, error)
}

// Files processes all files in the path that match the options
func Files(options *Options) (*Stats, error) {
	var handler = new(ioHandler)
	return processFiles(options, handler)
}

func processFiles(options *Options, ioHandler ioHandle) (*Stats, error) {

	data, err := ioHandler.ReadFile(options.LicensePath)
	if err != nil {
		return nil, err
	}
	license := string(data)

	channel := make(chan *Operation)
	start := time.Now()
	files := 0

	err = ioHandler.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
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

		data, err := ioHandler.ReadFile(path)
		if err != nil {
			files++
			channel <- &Operation{
				Action: OperationError,
				Path:   path,
			}
			return nil
		}

		fileContent := string(data)

		files++
		go func() {
			action := File(path, fileContent, license, options, ioHandler)
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
func File(filePath string, fileContent string, license string, options *Options, ioHandler ioHandle) Action {

	if strings.Contains(fileContent, strings.TrimSpace(license)) {
		return LicenseOk
	}

	if containsLicenseHeader(fileContent) {
		if options.Replace {
			newContent := replaceHeader(fileContent, license)
			if err := ioHandler.ReplaceFileContent(filePath, newContent); err != nil {
				return OperationError
			}
			return LicenseReplaced
		}
		return SkippedReplace
	}

	if options.Add {
		newContent := insertHeader(fileContent, license)
		if err := ioHandler.ReplaceFileContent(filePath, newContent); err != nil {
			return OperationError
		}
		return LicenseAdded
	}
	return SkippedAdd
}
