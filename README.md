![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Go/badge.svg)

# license-header-checker

![Demo](demo/demo.gif)

Command line utility written in [Go](https://golang.org) to **check** whether the **license headers** are included in all files of a specific project folder that match a specific list of extensions.

It can also **insert** the license in case it does not exist as well as **replace** wrong/old ones.

## Getting started

### Compiling

Provided that you have make installed, you can build & test the project by just typing:
```bash
$ make
```
If you do not want to use make, or you want more specific commands, take a look at the provided [Makefile](https://github.com/lsm-dev/license-header-checker/blob/master/Makefile).

### Installing

If this tool is going to be used frequently, it is often a good idea to install it in your path to be able to use it like any other command instead of using the executable's full path.

If you are new to Go, check that the [environment variables](https://golang.org/cmd/go/#hdr-Environment_variables) GOPATH and/or GOBIN are setup and added to your path).

To install the tool, just type:

```bash
$ make install
```

## How to use

### Usage

```bash
$ license-header-checker [option flags] /path/to/license_header /path/to/src [extensions...]
```

The option flags are:

```
-i 
	Insert the target license in case the file does not have any.

-r
	Replace the existing license by the target one in case it is 
	different (i.e. useful to change year)

-v
	Print extra information during execution like options, files 
	being processed, execution time, ...
````

### Example

This example assumes that:
* You have installed the tool (no need to use ./path/to/executable).
* You will execute the command directly from the folder you want to analyze.
* You have placed the "my_license_header.txt" file in the same folder.
* You are only interested in .js and .ts files.
* You want all options (verbose, insert, replace).

```bash
$ cd path/to/src/folder
$ license-header-checker -v -i -r ./my_license_header.txt . js ts
```
