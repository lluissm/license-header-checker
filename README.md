![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Go/badge.svg)

# license-header-checker

![Demo](demo/demo.gif)

Command line utility written in [Go](https://golang.org) to **check** whether the **license headers** are included in all files of a specific project folder that match a specific list of extensions.

It can also **insert** the license in case it does not exist as well as **replace** wrong/old ones.

## Compiling from source

Requires [go 1.13](https://golang.org/doc/devel/release.html#go1.13) and make. To build & test the project just type:

```bash
$ make
```
If you cannot use make or you want more specific commands, take a look at the provided [Makefile](https://github.com/lsm-dev/license-header-checker/blob/master/Makefile).

## Usage

Synopsis:

```bash
$ license-header-checker [-v] [-a] [-r] [-i folder1,folder2,...] <license-path> <project-path> <extensions...>
```

The following options are available:

```
  -a
    	Add the target license in case the file does not have any.

  -r
    	Replace the existing license by the target one in case they are different.

  -i
    	A comma separated list of the sub-folders that should be ignored.

  -v	
  		Be verbose during execution, printing options, files being processed, execution time, ...
````

Real-life example assuming that:
* You have installed the tool in your path (otherwise should start with ./path/to/executable).
* You will execute the command directly from the project folder (path to project => ".").
* You have placed the license_header.txt file in the same folder.
* You want all options (verbose, add, replace, ignore).
* You want to ignore node_modules and docs folders.
* You are only interested in .js and .ts files.


```bash
$ cd path/to/src/folder
$ license-header-checker -v -a -r -i node_modules,docs license_header.txt . js ts
```
