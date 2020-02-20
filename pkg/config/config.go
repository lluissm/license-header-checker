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

package config

import (
	"flag"
	"log"
	"strings"
)

// Options that the program accepts via commandline arguments/flags
type Options struct {
	Add         bool
	Replace     bool
	Verbose     bool
	Path        string
	LicensePath string
	Extensions  []string
	IgnorePaths []string
}

// ParseOptions returns the parsed options from cli
func ParseOptions() Options {

	writeFlagPtr := flag.Bool("a", false, "Add the target license in case the file does not have any.")
	overwriteFlagPtr := flag.Bool("r", false, "Replace the existing license by the target one in case they are different.")
	ignorePathsFlagPtr := flag.String("i", "", "A comma separated list of the sub-folders that should be ignored.")
	verboseFlagPtr := flag.Bool("v", false, "Be verbose during execution printing options, files being processed, execution time, ...")

	flag.Parse()
	args := flag.Args()

	if len(args) < 3 {
		log.Fatal("Missing arguments, please see documentation.")
	}

	licensePath := args[0]
	path := args[1]

	extensions := []string{}
	for _, e := range args[2:] {
		extensions = append(extensions, "."+e)
	}

	opt := &Options{
		*writeFlagPtr,
		*overwriteFlagPtr,
		*verboseFlagPtr,
		path,
		licensePath,
		extensions,
		strings.Split(*ignorePathsFlagPtr, ","),
	}

	return *opt
}
