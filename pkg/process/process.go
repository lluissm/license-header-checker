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
	"io/fs"
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
	// SkippedAdd means that the file had no license but the new one was not added.
	// as the -a flag was not provided
	SkippedAdd Action = iota
	// SkippedReplace means that the file had a different license but it was not
	// replaced with the target one as the -r flag was not provided
	SkippedReplace
	// LicenseOk means that the license was OK so the file was not changed
	LicenseOk
	// LicenseAdded means that the target license was added to the file
	LicenseAdded
	// LicenseReplaced means that the file's license was replaced by the target one
	LicenseReplaced
	// OperationError means there was an error with one of the files
	OperationError
)

// ioHandle defines the interface to manage files during processing
type ioHandle interface {
	// ReadFile returns the content of a file given its name
	ReadFile(name string) ([]byte, error)
	// WalkDir will call WalkDirFunc for every file in a directory
	WalkDir(path string, walkDirFn fs.WalkDirFunc) error
	// ReplaceFileContent replaces the content of the file with the one provided
	ReplaceFileContent(name string, content string) error
}

// File processes one file
func File(path string, content string, license string, options *Options, ioHandler ioHandle) Action {

	if strings.Contains(content, strings.TrimSpace(license)) {
		return LicenseOk
	}

	if containsLicenseHeader(content) {
		if options.Replace {
			newContent := replaceHeader(content, license)
			if err := ioHandler.ReplaceFileContent(path, newContent); err != nil {
				return OperationError
			}
			return LicenseReplaced
		}
		return SkippedReplace
	}

	if options.Add {
		newContent := insertHeader(content, license)
		if err := ioHandler.ReplaceFileContent(path, newContent); err != nil {
			return OperationError
		}
		return LicenseAdded
	}
	return SkippedAdd
}

// Files processes a group of files (in parallel) following the configuration
// defined in options
func Files(options *Options, ioHandler ioHandle) (*Stats, error) {

	data, err := ioHandler.ReadFile(options.LicensePath)
	if err != nil {
		return nil, err
	}

	license := string(data)
	channel := make(chan *Operation, 15)
	startTime := time.Now()
	stats := NewStats()
	files := 0

	err = ioHandler.WalkDir(options.Path, func(path string, d fs.DirEntry, err error) error {
		if processFile(channel, options, license, ioHandler, path, d, err) {
			files++
		}
		return nil
	})

	for i := 0; i < files; i++ {
		stats.AddOperation(<-channel)
	}

	stats.ElapsedMs = time.Since(startTime).Milliseconds()

	return stats, err
}

// processFile returns true if a file has been processed and false if processing has been skipped.
//
// Pocessing will be skipped if the path is a directory, it is part of the paths to ignore or the
// file extension does not match any of the extensions in options.Extensions
//
// The processing of the file is done on a goroutine, hence the channel to write the result of the
// operation
func processFile(channel chan *Operation, options *Options, license string, ioHandler ioHandle, path string, d fs.DirEntry, err error) bool {

	if d.IsDir() {
		return false
	}

	if shouldIgnorePath(path, options.IgnorePaths) {
		return false
	}

	if shouldIgnoreExtension(path, options.Extensions) {
		return false
	}

	if err != nil {
		onError(channel, path)
		return true
	}

	data, err := ioHandler.ReadFile(path)
	if err != nil {
		onError(channel, path)
		return true
	}

	go func() {
		content := string(data)
		action := File(path, content, license, options, ioHandler)
		channel <- &Operation{
			Action: action,
			Path:   path,
		}
	}()

	return true
}

func onError(channel chan *Operation, path string) {
	go func() {
		channel <- &Operation{
			Action: OperationError,
			Path:   path,
		}
	}()
}
