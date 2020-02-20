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
	"sync"
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
		os.Exit(0)
	}
	return string(data)
}

func processFiles(options config.Options) {
	files := 0
	start := time.Now()
	var wg sync.WaitGroup

	license := getLicense(options.LicensePath)

	licenseOk := 0
	licenseReplaced := 0
	licenseAdded := 0
	processedFiles := 0

	printIntro(options, license)

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
				res := processFile(path, license, options)

				switch res {
				case -1:
					licenseOk++
				case 0:
					licenseReplaced++
				case 1:
					licenseAdded++
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
	printResults(files, processedFiles, licenseOk, licenseAdded, licenseReplaced, elapsed.Milliseconds(), options)
}

func processFile(path string, license string, options config.Options) int {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content := string(data)

	if header.Contains(content, license) {
		return -1
	}

	if header.ContainsLicense(content) {
		if options.Verbose {
			fmt.Printf("%s: wrong license\n", path)
		}
		if options.Replace {
			replaceLicense(path, content, license)
		}
		return 0
	}

	if options.Verbose {
		fmt.Printf("%s: NO license\n", path)
	}
	if options.Add {
		insertLicense(path, content, license)
	}

	return 1
}

func insertLicense(path string, content string, license string) {
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

func printIntro(options config.Options, license string) {
	if options.Verbose {
		fmt.Printf("Project path: %s\n", options.Path)
		fmt.Printf("Ignore folders: %v\n", options.IgnorePaths)
		fmt.Printf("Extensions: %v\n", options.Extensions)
		fmt.Printf("Add license: %v\n", options.Add)
		fmt.Printf("Replace license: %v\n", options.Replace)
		fmt.Printf("Importing target license from: %s\n\n", options.LicensePath)
		fmt.Printf("%s\n\n", license)
	}
}

func printResults(files, processedFiles, licenseOk, licenseAdded, licenseReplaced int, elapsedMs int64, options config.Options) {
	red := color.FgRed.Render
	green := color.FgGreen.Render
	yellow := color.FgYellow.Render

	if options.Verbose {
		fmt.Printf("\nResults: %s ok, %s replaced, %s added, %d total files\n", green(fmt.Sprintf("%d", licenseOk)), yellow(fmt.Sprintf("%d", licenseReplaced)), red(fmt.Sprintf("%d", licenseAdded)), processedFiles)
	} else {
		fmt.Printf("Results: %s | %s | %s\n", green(fmt.Sprintf("%d", licenseOk)), yellow(fmt.Sprintf("%d", licenseReplaced)), red(fmt.Sprintf("%d", licenseAdded)))
	}

	if options.Verbose {
		fmt.Printf("\nTotal processing time: %d ms\n\n", elapsedMs)
	}
}
