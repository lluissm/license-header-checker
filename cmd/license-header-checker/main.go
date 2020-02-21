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

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gookit/color"
	"github.com/lsm-dev/license-header-checker/pkg/config"
	"github.com/lsm-dev/license-header-checker/pkg/file"
	"github.com/lsm-dev/license-header-checker/pkg/header"
)

func main() {
	options := config.ParseOptions()
	processFiles(options)
}

func getLicense(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

// Operation is the result of processFile
type Operation int

const (
	// Skipped means that the file was not modified but it did not have the target license
	Skipped Operation = iota
	// LicenseOk means that the license was OK
	LicenseOk
	// LicenseAdded means that the target license was added to the file
	LicenseAdded
	// LicenseReplaced means that the license was replaced with the target one
	LicenseReplaced
)

var (
	okRender      = color.FgGreen.Render
	warningRender = color.FgYellow.Render
	errorRender   = color.FgRed.Render
)

func processFiles(options config.Options) {
	start := time.Now()
	var wg sync.WaitGroup

	license := getLicense(options.LicensePath)

	var files, processedFiles, licenseOk, licenseAdded, licenseReplaced, skipped int64 = 0, 0, 0, 0, 0, 0

	printIntro(options)

	err := filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if file.ShouldIgnore(path, options.IgnorePaths) {
			return nil
		}

		if file.HasExtension(path, options.Extensions) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				op, _ := processFile(path, license, options)

				switch op {
				case Skipped:
					atomic.AddInt64(&skipped, 1)
				case LicenseOk:
					atomic.AddInt64(&licenseOk, 1)
				case LicenseReplaced:
					atomic.AddInt64(&licenseReplaced, 1)
				case LicenseAdded:
					atomic.AddInt64(&licenseAdded, 1)
				}
			}()
			processedFiles++
		}

		files++
		return nil
	})

	if err != nil {
		log.Println(err)
	}

	wg.Wait()

	elapsed := time.Since(start)
	printResults(files, processedFiles, skipped, licenseOk, licenseAdded, licenseReplaced, elapsed.Milliseconds(), options)
}

func processFile(path string, license string, options config.Options) (Operation, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("%s\n", errorRender("%s", err))
		return Skipped, err
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
			replaceLicense(path, content, license)
			return LicenseReplaced, nil
		}
		return Skipped, nil
	}

	if options.Verbose {
		fmt.Printf(" · %s => %s", path, errorRender("License missing\n"))
	}
	if options.Add {
		addLicense(path, content, license)
		return LicenseAdded, nil
	}

	return Skipped, nil
}

func addLicense(path string, content string, license string) {
	res := header.Insert(content, license)
	replaceFile(path, res)
}

func replaceLicense(path string, content string, license string) {
	res := header.Replace(content, license)
	replaceFile(path, res)
}

func replaceFile(path string, content string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("failed deleting the file: %s", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	if err = writer.Flush(); err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func printIntro(options config.Options) {
	if options.Verbose {
		blue := color.FgBlue.Render
		fmt.Printf("Options: ")
		fmt.Printf("\n · Project path: %s\n", blue(fmt.Sprintf("%s", options.Path)))
		fmt.Printf(" · Ignore folders: %s\n", blue(fmt.Sprintf("%v", options.IgnorePaths)))
		fmt.Printf(" · Extensions: %s\n", blue(fmt.Sprintf("%v", options.Extensions)))
		fmt.Printf(" · Flags: ")
		if options.Add {
			fmt.Printf("%s ", blue("add"))
		}
		if options.Replace {
			fmt.Printf("%s ", blue("replace"))
		}
		fmt.Printf("\n · License header: %s\n", blue(fmt.Sprintf("%s", options.LicensePath)))
		fmt.Printf("\nFiles:\n")
	}
}

func printResults(files, processedFiles, skipped, licensesOk, licensesAdded, licensesReplaced, elapsedMs int64, options config.Options) {
	licensesOkStr := okRender(fmt.Sprintf("%d", licensesOk))
	licensesReplacedStr := warningRender(fmt.Sprintf("%d", licensesReplaced))
	licensesAddedStr := errorRender(fmt.Sprintf("%d", licensesAdded))

	if options.Verbose {
		fmt.Printf("\nResults:\n")
		fmt.Printf(" · Total files: %d\n", processedFiles)
		fmt.Printf(" · OK licenses: %s\n", licensesOkStr)
		fmt.Printf(" · Added licenses: %s\n", licensesAddedStr)
		fmt.Printf(" · Replaced licenses: %s\n", licensesReplacedStr)
		fmt.Printf(" · Processing time: %d ms\n", elapsedMs)
	} else {
		fmt.Printf("%d files => %s licenses ok, %s licenses replaced, %s licenses added", processedFiles, licensesOkStr, licensesReplacedStr, licensesAddedStr)
	}

	if skipped > 0 {
		fmt.Println("")
		color.Error.Printf("\n [!] %d operations were skipped. You may have forgotten to add one of the following flags -a (add), -r (replace) ", skipped)
	}
}
