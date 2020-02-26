![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Go/badge.svg)

# license-header-checker

![Demo](demo/demo.gif)

Command line utility written in [Go](https://golang.org) to **check** whether the **license headers** are included in the **source files** of a project.

Optionally, the tool can **insert** the license in case it does not exist as well as **replace** it in case it is obsolete.

## Usage

Syntax:

```bash
$ license-header-checker [-a] [-r] [-v] [-i dir1,dir2,file1,path1...] license-header-path project-path extensions...
```

Options:

```
  -a        Add the target license in case the file does not have any.
  -r        Replace the existing license by the target one in case they are different.
  -v        Be verbose during execution.
  -i        A comma separated list of the folders/files/paths that should be ignored. 
	        It does not support wildcards.
  -version  Display version number.
```

Example:

```bash
$ license-header-checker -v -a -r -i node_modules,myFile.js,client/assets ../license_header.txt . js ts
```

## Installation


### Binary packages

The binary packages for Linux, Windows and macOS are uploaded for each [release](https://github.com/lsm-dev/license-header-checker/releases).



### Building from source

Requires [go 1.13](https://golang.org/doc/devel/release.html#go1.13) and `make`.

To **build**:

```bash
$ make build
```

To execute **unit tests**:

```bash
$ make test
```

To **install** in go/bin:

```bash
$ make install
```

To **cross-compile** (generate the binaries for Linux, Windows and macOS):

```bash
$ make cross-build
```
