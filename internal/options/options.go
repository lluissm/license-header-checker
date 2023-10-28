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

package options

import (
	"errors"
	"flag"
	"fmt"
	"github.com/lluissm/license-header-checker/pkg/process"
	"regexp"
	"strings"
)

// Options are the process.Options parsed from command line flags/args
type Options struct {
	ShowVersion bool
	Verbose     bool
	Process     *process.Options
}

// Parse returns the parsed Options from command line flags/args
func Parse(osArgs []string) (*Options, error) {

	flagSet := flag.NewFlagSet("lhc", flag.ExitOnError)
	flagSet.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\033[1;4mSYNOPSIS\033[0m\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "license-header-checker [-a] [-r] [-v] [-i path1,...] license-header-path src-path extensions...\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\033[1;4mOPTIONS\033[0m\n\n")
		flagSet.PrintDefaults()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\n\033[1;4mEXAMPLE\033[0m\n\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "license-header-checker -a -r -v -i folder,ignore/path license-header-path project-src-path extension1 extension2\n\n")
	}

	addFlag := flagSet.Bool("a", false, "Add the target license in case the file does not have any.")
	replaceFlag := flagSet.Bool("r", false, "Replace the existing license by the target one in case they are different.")
	ignorePathsFlag := flagSet.String("i", "", "A comma separated list of the folders, files and/or paths that should be ignored. Does not support wildcards.")
	verboseFlag := flagSet.Bool("v", false, "Be verbose during execution printing options, files being processed, execution time, ...")
	headerRegexFlag := flagSet.String("e", "", "A regular expression to match a header comment.")
	showVersionFlag := flagSet.Bool("version", false, "Display version number")

	if err := flagSet.Parse(osArgs[1:]); err != nil {
		return nil, err
	}

	args := flagSet.Args()

	if *showVersionFlag {
		return &Options{
			true,
			false,
			nil,
		}, nil
	}

	if len(args) < 3 {
		return nil, errors.New("missing arguments, please see documentation")
	}

	licensePath := args[0]
	path := args[1]

	var extensions []string
	for _, e := range args[2:] {
		extensions = append(extensions, "."+e)
	}

	var ignorePaths []string
	for _, p := range strings.Split(*ignorePathsFlag, ",") {
		if len(p) > 0 {
			ignorePaths = append(ignorePaths, p)
		}
	}

	var headerRegex = process.DefaultRegex
	if headerRegexFlag != nil && len(*headerRegexFlag) > 0 {
		rex, err := regexp.Compile(*headerRegexFlag)
		if err != nil {
			return nil, err
		}
		headerRegex = rex
	}

	processOptions := &process.Options{
		Add:         *addFlag,
		Replace:     *replaceFlag,
		Path:        path,
		LicensePath: licensePath,
		Extensions:  extensions,
		IgnorePaths: ignorePaths,
		HeaderRegex: headerRegex,
	}

	return &Options{
		*showVersionFlag,
		*verboseFlag,
		processOptions,
	}, nil
}
