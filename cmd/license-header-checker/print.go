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
		printFiles(stats)
		printOptions(options)
		printTotals(stats)
	} else {
		printShort(stats)
	}
	printWarnings(stats)
}

// printShort prints the result of the processing in a compact mode (non-verbose)
func printShort(stats *process.Stats) {
	fmt.Printf("%s licenses ok, %s licenses replaced, %s licenses added\n",
		okRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseOk]))),
		warningRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseReplaced]))),
		errorRender(fmt.Sprintf("%d", len(stats.Files[process.LicenseAdded]))))
}

// printFiles prints the the output of each file grouped by the result type
func printFiles(stats *process.Stats) {
	fmt.Printf("files:\n")
	if licensesOk := stats.Files[process.LicenseOk]; len(licensesOk) > 0 {
		fmt.Printf("  license_ok:\n")
		for _, file := range licensesOk {
			fmt.Printf("    - %s\n", okRender(fmt.Sprintf("%v", file)))
		}
	}
	if licensesReplaced := stats.Files[process.LicenseReplaced]; len(licensesReplaced) > 0 {
		fmt.Printf("  license_replaced:\n")
		for _, file := range licensesReplaced {
			fmt.Printf("    - %s\n", warningRender(fmt.Sprintf("%v", file)))
		}
	}
	if licensesAdded := stats.Files[process.LicenseAdded]; len(licensesAdded) > 0 {
		fmt.Printf("  license_added:\n")
		for _, file := range licensesAdded {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if skippedAdds := stats.Files[process.SkippedAdd]; len(skippedAdds) > 0 {
		fmt.Printf("  skipped_add:\n")
		for _, file := range skippedAdds {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if skippedReplaces := stats.Files[process.SkippedReplace]; len(skippedReplaces) > 0 {
		fmt.Printf("  skipped_replace:\n")
		for _, file := range skippedReplaces {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if errors := stats.Files[process.OperationError]; len(errors) > 0 {
		fmt.Printf("  errors:\n")
		for _, file := range errors {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
}

// printOptions prints the options that were supplied to the cli app
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

// printTotals prints the aggregated data
func printTotals(stats *process.Stats) {
	fmt.Printf("totals:\n")
	if licensesOk := len(stats.Files[process.LicenseOk]); licensesOk > 0 {
		fmt.Printf("  license_ok: %s\n", okRender(fmt.Sprintf("%v files", licensesOk)))
	}
	if licensesReplaced := len(stats.Files[process.LicenseReplaced]); licensesReplaced > 0 {
		fmt.Printf("  license_replaced: %s\n", warningRender(fmt.Sprintf("%v files", licensesReplaced)))
	}
	if licensesAdded := len(stats.Files[process.LicenseAdded]); licensesAdded > 0 {
		fmt.Printf("  license_added: %s\n", errorRender(fmt.Sprintf("%v files", licensesAdded)))
	}
	if skippedAdds := len(stats.Files[process.SkippedAdd]); skippedAdds > 0 {
		fmt.Printf("  skipped_add: %s\n", errorRender(fmt.Sprintf("%v files", skippedAdds)))
	}
	if skippedReplaces := len(stats.Files[process.SkippedReplace]); skippedReplaces > 0 {
		fmt.Printf("  skipped_replace: %s\n", errorRender(fmt.Sprintf("%v files", skippedReplaces)))
	}
	if errors := len(stats.Files[process.OperationError]); errors > 0 {
		fmt.Printf("  error: %s\n", errorRender(fmt.Sprintf("%v files", errors)))
	}
	fmt.Printf("  elapsed_time: %s\n", infoRender(fmt.Sprintf("%vms", stats.ElapsedMs)))
}

// printWarnings warns the user that the -a or -r flag were not provided
// and they may have had been useful
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
