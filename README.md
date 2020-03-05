![Header](images/header.png)

![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Build/badge.svg)   ![Test Status](https://github.com/lsm-dev/license-header-checker/workflows/Test/badge.svg)

Multiplatform command line tool that **checks** whether the **license headers** are included in the **source files** of a project.

- It can **insert the license** to a file in case it does not exist (optional).
- It can **replace the license** of a file in case it is different (optional).
- It is possible to **choose the file extensions that will be processed**.
- Specific **folders, filers or paths can be ignored**.

_DISCLAIMER_

The tool expects the software license to be in a **block comment at the beginning of the file** following the format `/*  */`. 

Thus, while it does support the source files of languages like *Go, Rust, JavaScript, TypeScript, C, C++, Java, Swift, Kotlin and C#*, it does not support the file extensions that do not use this style.

# Usage

## Syntax

```bash
$ license-header-checker [-a] [-r] [-v] [-i path1,...] license-header-path src-path extensions...
```

## Options

```
  -a        Add the target license in case the file does not have any.
  -r        Replace the existing license by the target one in case they are different.
  -v        Be verbose during execution.
  -i        A comma separated list of the folders, files and/or paths that should be ignored. 
            It does not support wildcards.
  -version  Display version number.
```

## Example

```bash
$ license-header-checker -v -a -r -i node_modules,client/assets ../license_header.txt . js ts
```

# Installation

## Binary packages

The binary packages for Linux, Windows and macOS are uploaded for each release and can be downloaded from [here](https://github.com/lsm-dev/license-header-checker/releases).

## Building from source
The tool has been built with [go 1.13](https://golang.org/doc/devel/release.html#go1.13).

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
