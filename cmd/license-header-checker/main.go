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
	"log"
	"os"

	"github.com/lluissm/license-header-checker/internal/options"
	"github.com/lluissm/license-header-checker/pkg/process"
)

var version string = "development"

func main() {

	options, err := options.Parse(os.Args)
	if err != nil {
		log.Fatalf("could not parse the cli args: %s", err.Error())
	}

	if options.ShowVersion {
		fmt.Printf("version: %s\n", version)
		os.Exit(0)
	}

	stats, err := process.Files(options.Process, new(ioHandler))
	if err != nil {
		log.Fatalf("could not process the files: %s", err.Error())
	}

	printStats(options, stats)
}
