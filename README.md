![Demo](demo/demo.gif)

# license-header-checker

![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Build/badge.svg)   ![Test Status](https://github.com/lsm-dev/license-header-checker/workflows/Test/badge.svg)

Command line utility written in [Go](https://golang.org) to **check** whether the **license headers** are included in the **source files** of a project.

Optionally, the tool can **insert** the license to a file in case it does not exist as well as **replace** it in case it is obsolete.

The recommendation is to use the tool directly in the root folder of the project as only the files that match the provided **extensions** will be analyzed. Furthermore, there is the possibility to **ignore specific files and folders** if necessary by providing the `-i` flag.

The tool expects the software license to be in a **block comment at the beginning of the file** following the format `/*  */`. 

Thus, while it does support the source files of languages like *Go, Rust, JavaScript, TypeScript, C, C++, Java, Swift, Kotlin and C#*, it does not support the file extensions that do not use this style.

The tool has been built with [go 1.13](https://golang.org/doc/devel/release.html#go1.13) and the build process is automated with `make`.

## Usage

Syntax:

```bash
$ license-header-checker [-a] [-r] [-v] [-i path1,...] license-header-path src-path extensions...
```

Options:

```
  -a        Add the target license in case the file does not have any.
  -r        Replace the existing license by the target one in case they are different.
  -v        Be verbose during execution.
  -i        A comma separated list of the folders, files and/or paths that should be ignored. 
            It does not support wildcards.
  -version  Display version number.
```

Example:

```bash
$ license-header-checker -v -a -r -i node_modules,client/assets ../license_header.txt . js ts
```

## Installation

### Binary packages

The binary packages for Linux, Windows and macOS are uploaded for each release and can be downloaded from [here](https://github.com/lsm-dev/license-header-checker/releases).

### Building from source


To **build**:

```bash
$ make build
```

To **install** in go/bin:

```bash
$ make install
```

To **cross-compile** (generate the binaries for Linux, Windows and macOS all at once):

```bash
$ make cross-build
```
To run **unit** tests:

```bash
$ make test
```
To run **end-to-end** tests:

```bash
$ make test-e2e
```
