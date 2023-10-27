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

package main

import (
	"fmt"
	"sort"

	"github.com/gookit/color"
	"github.com/lluissm/license-header-checker/internal/options"
	"github.com/lluissm/license-header-checker/pkg/process"
)

var (
	okRender      = color.FgGreen.Render
	infoRender    = color.FgBlue.Render
	warningRender = color.FgYellow.Render
	errorRender   = color.FgRed.Render
)

// printStats writes to the standard output the result of the processing according
// to the verbosity level
func printStats(options *options.Options, stats *process.Stats) {
	if options.Verbose {
		printFileOperations(stats)
		printOptions(options)
		printTotals(stats)
	} else {
		printShort(stats)
	}
	printWarnings(stats)
}

// printFileOperations prints the files processed by operation type
func printFileOperations(stats *process.Stats) {
	fmt.Printf("files:\n")
	printFiles(stats.Files[process.LicenseOk], "license_ok", okRender)
	printFiles(stats.Files[process.LicenseReplaced], "license_replaced", warningRender)
	printFiles(stats.Files[process.LicenseAdded], "license_added", errorRender)
	printFiles(stats.Files[process.SkippedAdd], "skipped_add", errorRender)
	printFiles(stats.Files[process.SkippedReplace], "skipped_replace", errorRender)
	printFiles(stats.Files[process.OperationError], "errors", errorRender)
}

// printOptions prints the options that were supplied to the app
func printOptions(options *options.Options) {
	fmt.Printf("options:\n")
	fmt.Printf("  project_path: %s\n", infoRender(options.Process.Path))
	if len(options.Process.IgnorePaths) > 0 {
		fmt.Printf("  ignore_paths:\n")
		for _, ignorePaths := range options.Process.IgnorePaths {
			fmt.Printf("    - %s\n", infoRender(fmt.Sprintf("%v", ignorePaths)))
		}
	}
	fmt.Printf("  extensions:\n")
	for _, ext := range options.Process.Extensions {
		fmt.Printf("    - %s\n", infoRender(fmt.Sprintf("%v", ext)))
	}
	fmt.Printf("  flags:\n")
	if options.Process.Add {
		fmt.Printf("    - %s\n", infoRender("add"))
	}
	if options.Process.Replace {
		fmt.Printf("    - %s\n", infoRender("replace"))
	}
	if options.Verbose {
		fmt.Printf("    - %s\n", infoRender("verbose"))
	}
	fmt.Printf("  license_header: %s\n", infoRender("%s", options.Process.LicensePath))
}

// printTotals prints the total amount of files processed by operation type
func printTotals(stats *process.Stats) {
	fmt.Printf("totals:\n")
	printFileTotals(len(stats.Files[process.LicenseOk]), "license_ok", okRender)
	printFileTotals(len(stats.Files[process.LicenseReplaced]), "license_replaced", warningRender)
	printFileTotals(len(stats.Files[process.LicenseAdded]), "license_added", errorRender)
	printFileTotals(len(stats.Files[process.SkippedAdd]), "skipped_add", errorRender)
	printFileTotals(len(stats.Files[process.SkippedReplace]), "skipped_replace", errorRender)
	printFileTotals(len(stats.Files[process.OperationError]), "error", errorRender)
	fmt.Printf("  elapsed_time: %s\n", infoRender(fmt.Sprintf("%vms", stats.ElapsedMs)))
}

// printShort prints the result of the processing in a compact mode (non-verbose)
func printShort(stats *process.Stats) {
	fmt.Printf("%s licenses ok, %s licenses replaced, %s licenses added\n",
		okRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseOk]))),
		warningRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseReplaced]))),
		errorRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseAdded]))))
}

// printWarnings warns the user if the -a or -r flag were not provided
// but they may have had been useful
func printWarnings(stats *process.Stats) {
	if skippedAdds := len(stats.Files[process.SkippedAdd]); skippedAdds > 0 {
		color.Error.Printf("[!] %d files had no license but were not changed as the -a (add) option was not supplied.\n", skippedAdds)
	}
	if skippedReplaces := len(stats.Files[process.SkippedReplace]); skippedReplaces > 0 {
		color.Error.Printf("[!] %d files had a different license but were not changed as the -r (replace) option was not supplied.\n", skippedReplaces)
	}
	if errors := len(stats.Files[process.OperationError]); errors > 0 {
		color.Error.Printf("[!] There where %d errors.\n", errors)
	}
}

func printFiles(files []string, operationName string, render func(a ...interface{}) string) {
	if len(files) <= 0 {
		return
	}
	fmt.Printf("  %s:\n", operationName)
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	for _, file := range files {
		fmt.Printf("    - %s\n", render(fmt.Sprintf("%v", file)))
	}
}

func printFileTotals(total int, operationName string, render func(a ...interface{}) string) {
	if total <= 0 {
		return
	}
	fmt.Printf("  %s: %s\n", operationName, render(fmt.Sprintf("%v files", total)))
}
