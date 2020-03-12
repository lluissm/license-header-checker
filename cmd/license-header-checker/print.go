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
	"fmt"

	"github.com/gookit/color"
	"github.com/lsm-dev/license-header-checker/pkg/process"
)

var (
	okRender      = color.FgGreen.Render
	infoRender    = color.FgBlue.Render
	warningRender = color.FgYellow.Render
	errorRender   = color.FgRed.Render
)

// printStats prints the options, files processed and aggregate data
func printStats(verbose bool, stats *process.Stats, options *process.Options) {
	var licensesOk, licensesAdded, licensesReplaced, skippedAdds, skippedReplaces, errors []string

	for _, op := range stats.Operations {
		action := op.Action
		switch action {
		case process.SkippedAdd:
			skippedAdds = append(skippedAdds, op.Path)
		case process.SkippedReplace:
			skippedReplaces = append(skippedReplaces, op.Path)
		case process.LicenseOk:
			licensesOk = append(licensesOk, op.Path)
		case process.LicenseAdded:
			licensesAdded = append(licensesAdded, op.Path)
		case process.LicenseReplaced:
			licensesReplaced = append(licensesReplaced, op.Path)
		case process.OperationError:
			errors = append(errors, op.Path)
		}
	}

	if verbose {
		printProcessedFiles(licensesOk, licensesReplaced, licensesAdded, skippedReplaces, skippedAdds, errors)
		printProcessOptions(options)
		printTotals(len(licensesOk), len(licensesReplaced), len(licensesAdded), len(skippedReplaces), len(skippedAdds), len(errors), int(stats.ElapsedMs))
	} else {
		fmt.Printf("%s licenses ok, %s licenses replaced, %s licenses added\n", okRender(fmt.Sprintf("%d", len(licensesOk))), warningRender(fmt.Sprintf("%d", len(licensesReplaced))), errorRender(fmt.Sprintf("%d", len(licensesAdded))))
	}

	printWarnings(len(skippedReplaces), len(skippedAdds), len(errors))
}

// printProcessedFiles prints the the output of each file grouped by result type
func printProcessedFiles(licensesOk, licensesReplaced, licensesAdded, skippedReplaces, skippedAdds, errors []string) {
	fmt.Printf("files:\n")
	if len(licensesOk) > 0 {
		fmt.Printf("  license_ok:\n")
		for _, file := range licensesOk {
			fmt.Printf("    - %s\n", okRender(fmt.Sprintf("%v", file)))
		}
	}
	if len(licensesReplaced) > 0 {
		fmt.Printf("  license_replaced:\n")
		for _, file := range licensesReplaced {
			fmt.Printf("    - %s\n", warningRender(fmt.Sprintf("%v", file)))
		}
	}
	if len(licensesAdded) > 0 {
		fmt.Printf("  license_added:\n")
		for _, file := range licensesAdded {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if len(skippedAdds) > 0 {
		fmt.Printf("  skipped_add:\n")
		for _, file := range skippedAdds {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if len(skippedReplaces) > 0 {
		fmt.Printf("  skipped_replace:\n")
		for _, file := range skippedReplaces {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
	if len(errors) > 0 {
		fmt.Printf("  errors:\n")
		for _, file := range errors {
			fmt.Printf("    - %s\n", errorRender(fmt.Sprintf("%v", file)))
		}
	}
}

// printProcessOptions prints the options that were supplied to the cli app
func printProcessOptions(options *process.Options) {
	fmt.Printf("options:\n")
	fmt.Printf("  project_path: %s\n", infoRender(fmt.Sprintf("%s", options.Path)))
	if len(options.IgnorePaths) > 0 {
		fmt.Printf("  ignore_paths:\n")
		for _, ignorePaths := range options.IgnorePaths {
			fmt.Printf("    - %s\n", infoRender(fmt.Sprintf("%v", ignorePaths)))
		}
	}
	fmt.Printf("  extensions:\n")
	for _, ext := range options.Extensions {
		fmt.Printf("    - %s\n", infoRender(fmt.Sprintf("%v", ext)))
	}
	fmt.Printf("  flags:\n")
	if options.Add {
		fmt.Printf("    - %s\n", infoRender("add"))
	}
	if options.Replace {
		fmt.Printf("    - %s\n", infoRender("replace"))
	}
	fmt.Printf("  license_header: %s\n", infoRender(fmt.Sprintf("%s", options.LicensePath)))
}

// printTotals prints the aggregated data
func printTotals(licensesOk, licensesReplaced, licensesAdded, skippedReplaces, skippedAdds, errors, elapsedMs int) {
	fmt.Printf("totals:\n")
	if licensesOk > 0 {
		fmt.Printf("  license_ok: %s\n", okRender(fmt.Sprintf("%v files", licensesOk)))
	}
	if licensesReplaced > 0 {
		fmt.Printf("  license_replaced: %s\n", warningRender(fmt.Sprintf("%v files", licensesReplaced)))
	}
	if licensesAdded > 0 {
		fmt.Printf("  license_added: %s\n", errorRender(fmt.Sprintf("%v files", licensesAdded)))
	}
	if skippedAdds > 0 {
		fmt.Printf("  skipped_add: %s\n", errorRender(fmt.Sprintf("%v files", skippedAdds)))
	}
	if skippedReplaces > 0 {
		fmt.Printf("  skipped_replace: %s\n", errorRender(fmt.Sprintf("%v files", skippedReplaces)))
	}
	if errors > 0 {
		fmt.Printf("  error: %s\n", errorRender(fmt.Sprintf("%v files", errors)))
	}
	fmt.Printf("  elapsed_time: %s\n", infoRender(fmt.Sprintf("%vms", elapsedMs)))
}

// printWarnings prints warnings to the user in case he/she may have forgotten to add -a or -r flag
func printWarnings(skippedReplaces, skippedAdds, errors int) {
	if skippedAdds > 0 {
		color.Error.Printf("[!] %d files had no license but were not changed as the -a (add) option was not supplied.\n", skippedAdds)
	}
	if skippedReplaces > 0 {
		color.Error.Printf("[!] %d files had a different license but were not changed as the -r (replace) option was not supplied.\n", skippedReplaces)
	}
	if errors > 0 {
		color.Error.Printf("[!] There where %d errors.\n", errors)
	}
}
