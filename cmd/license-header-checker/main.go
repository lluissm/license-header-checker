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
		log.Fatal(err.Error())
	}
	if options.ShowVersion {
		fmt.Printf("version: %s\n", version)
		os.Exit(0)
	}
	stats, err := process.Files(options.Process, new(ioHandler))
	if err != nil {
		log.Fatal(err.Error())
	}
	print(options, stats)
}
