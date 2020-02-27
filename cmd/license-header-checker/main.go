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
	"os"

	"github.com/gookit/color"
	"github.com/lsm-dev/license-header-checker/internal/config"
	"github.com/lsm-dev/license-header-checker/internal/process"
)

func main() {
	options := config.ParseOptions()
	stats, err := process.Files(options)
	if err != nil {
		fmt.Println(errorRender(err))
		os.Exit(1)
	}
	printOptions(options)
	printStats(stats, options)
}

var (
	okRender      = color.FgGreen.Render
	infoRender    = color.FgBlue.Render
	warningRender = color.FgYellow.Render
	errorRender   = color.FgRed.Render
)

func printOptions(options *config.Options) {
	if options.Verbose {

		fmt.Printf("\nOptions: ")
		fmt.Printf("\n · Project path: %s\n", infoRender(fmt.Sprintf("%s", options.Path)))
		fmt.Printf(" · Ignore folders: %s\n", infoRender(fmt.Sprintf("%v", options.IgnorePaths)))
		fmt.Printf(" · Extensions: %s\n", infoRender(fmt.Sprintf("%v", options.Extensions)))
		fmt.Printf(" · Flags: ")
		if options.Add {
			fmt.Printf("%s ", infoRender("add"))
		}
		if options.Replace {
			fmt.Printf("%s ", infoRender("replace"))
		}
		fmt.Printf("\n · License header: %s\n", infoRender(fmt.Sprintf("%s", options.LicensePath)))
		fmt.Printf("\nFiles:\n")
	}
}

func printStats(stats *process.Stats, options *config.Options) {
	licensesOkStr := okRender(fmt.Sprintf("%d", stats.LicensesOk))
	licensesReplacedStr := warningRender(fmt.Sprintf("%d", stats.LicensesReplaced))
	licensesAddedStr := errorRender(fmt.Sprintf("%d", stats.LicensesAdded))

	totalFiles := stats.TotalOperations()

	if options.Verbose {
		fmt.Printf("\nResults:\n")
		fmt.Printf(" · Total files: %d\n", totalFiles)
		fmt.Printf(" · OK licenses: %s\n", licensesOkStr)
		fmt.Printf(" · Added licenses: %s\n", licensesAddedStr)
		fmt.Printf(" · Replaced licenses: %s\n", licensesReplacedStr)
		fmt.Printf(" · Processing time: %d ms\n", stats.ElapsedMs)
	} else {
		fmt.Printf("%d files => %s licenses ok, %s licenses replaced, %s licenses added\n", totalFiles, licensesOkStr, licensesReplacedStr, licensesAddedStr)
	}

	if stats.SkippedAdds > 0 {
		color.Error.Printf("[!] %d files had no license but where not changed as the -a (add) option was not supplied.\n", stats.SkippedAdds)
	}
	if stats.SkippedReplaces > 0 {
		color.Error.Printf("[!] %d files had a different license but where not changed as the -r (replace) option was not supplied.\n", stats.SkippedReplaces)
	}
	if stats.Errors > 0 {
		color.Error.Printf("[!] There where %d errors.\n", stats.Errors)
	}
}
