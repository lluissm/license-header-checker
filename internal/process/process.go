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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/lsm-dev/license-header-checker/internal/config"
	"github.com/lsm-dev/license-header-checker/internal/file"
	"github.com/lsm-dev/license-header-checker/internal/header"
)

func getLicense(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type (
	// Operation is the result of processing one file
	Operation int

	// Stats is the result of processing multiple files
	Stats struct {
		SkippedAdds      int64
		SkippedReplaces  int64
		LicensesOk       int64
		LicensesAdded    int64
		LicensesReplaced int64
		Errors           int64
		ElapsedMs        int64
	}
)

// TotalOperations returns the sum of Skipped, LicensesOk, LicensesAdded and LicensesReplaced
func (s *Stats) TotalOperations() int64 {
	return s.SkippedAdds + s.SkippedReplaces + s.LicensesOk + s.LicensesAdded + s.LicensesReplaced
}

const (
	// SkippedAdd means that the file had no license but the new one was not added. Missing -a flasg
	SkippedAdd Operation = iota
	// SkippedReplace means that the file had a different license but it was not replaced with the target one. Missing -r flag
	SkippedReplace
	// LicenseOk means that the license was OK
	LicenseOk
	// LicenseAdded means that the target license was added to the file
	LicenseAdded
	// LicenseReplaced means that the license was replaced with the target one
	LicenseReplaced
	// OperationError means there was an error with one of the files
	OperationError
)

var (
	okRender      = color.FgGreen.Render
	warningRender = color.FgYellow.Render
	errorRender   = color.FgRed.Render
)

// Files processes all files matching with options
func Files(options *config.Options) (*Stats, error) {
	start := time.Now()
	var wg sync.WaitGroup

	license, err := getLicense(options.LicensePath)
	if err != nil {
		return nil, err
	}

	var licenseOk, licenseAdded, licenseReplaced, skippedAdds, skippedReplaces, errors int64 = 0, 0, 0, 0, 0, 0

	err = filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			fmt.Printf("%s\n", errorRender("%s", err))
			atomic.AddInt64(&errors, 1)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if file.ShouldIgnore(path, options.IgnorePaths) {
			return nil
		}
		if !file.HasExtension(path, options.Extensions) {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			op, err := File(path, license, options)
			if err != nil {
				fmt.Printf("%s\n", errorRender("%s", err))
			}

			switch op {
			case SkippedAdd:
				atomic.AddInt64(&skippedAdds, 1)
			case SkippedReplace:
				atomic.AddInt64(&skippedReplaces, 1)
			case LicenseOk:
				atomic.AddInt64(&licenseOk, 1)
			case LicenseReplaced:
				atomic.AddInt64(&licenseReplaced, 1)
			case LicenseAdded:
				atomic.AddInt64(&licenseAdded, 1)
			case OperationError:
				atomic.AddInt64(&errors, 1)
			}
		}()

		return nil
	})

	wg.Wait()

	elapsed := time.Since(start)
	return &Stats{
		SkippedAdds:      skippedAdds,
		SkippedReplaces:  skippedReplaces,
		LicensesOk:       licenseOk,
		LicensesAdded:    licenseAdded,
		LicensesReplaced: licenseReplaced,
		ElapsedMs:        elapsed.Milliseconds(),
		Errors:           errors,
	}, err
}

// File processes the file
func File(path string, license string, options *config.Options) (Operation, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return OperationError, err
	}

	content := string(data)

	if strings.Contains(content, license) {
		if options.Verbose {
			fmt.Printf(" · %s => %s", path, okRender("License ok\n"))
		}
		return LicenseOk, nil
	}

	if header.ContainsLicense(content) {
		if options.Verbose {
			fmt.Printf(" · %s => %s", path, warningRender("License is different\n"))
		}
		if options.Replace {
			if err := replaceLicense(path, content, license); err != nil {
				return OperationError, err
			}
			return LicenseReplaced, nil
		}
		return SkippedReplace, nil
	}

	if options.Verbose {
		fmt.Printf(" · %s => %s", path, errorRender("License missing\n"))
	}
	if options.Add {
		if err := addLicense(path, content, license); err != nil {
			return OperationError, err
		}
		return LicenseAdded, nil
	}
	return SkippedAdd, nil
}

func addLicense(path string, content string, license string) error {
	res := header.Insert(content, license)
	return file.Replace(path, res)
}

func replaceLicense(path string, content string, license string) error {
	res := header.Replace(content, license)
	return file.Replace(path, res)
}
