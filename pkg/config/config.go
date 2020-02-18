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
)

// Options that the program accepts via commandline arguments/flags
type Options struct {
	Insert      bool
	Replace     bool
	Verbose     bool
	Path        string
	LicensePath string
	Extensions  []string
}

// ParseOptions returns the options
func ParseOptions() Options {

	writeFlagPtr := flag.Bool("i", false, "Insert the target license in case the file does not have any.")
	overwriteFlagPtr := flag.Bool("r", false, "Replace the existing license by the target one in case it is different (i.e. useful to change year)")
	verboseFlagPtr := flag.Bool("v", false, "Print extra information during execution like options, files being processed, execution time, ...")

	flag.Parse()

	args := flag.Args()
	if len(args) < 3 {
		panic("Missing argument")
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
	}
	return *opt
}
