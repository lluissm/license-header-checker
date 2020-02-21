![Build Status](https://github.com/lsm-dev/license-header-checker/workflows/Go/badge.svg)

# license-header-checker

![Demo](demo/demo.gif)

Command line utility written in [Go](https://golang.org) to **check** whether the **license headers** are included in all files of a specific project folder that match a specific list of extensions.

It can also **insert** the license in case it does not exist as well as **replace** wrong/old ones.

## Compiling from source

Requires [go 1.13](https://golang.org/doc/devel/release.html#go1.13) and make. To **build & test** the project:

```bash
$ make
```



To generate the **binaries for Windows, mac OS and linux**:

```bash
$ make cross-build
```

## Usage

Synopsis:

```bash
$ license-header-checker [-a] [-r] [-v] [-i folder1,folder2,...] license-path project-path extensions...
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
	Be verbose during execution printing options, files being processed, 
	execution time, ...
````

Example:

```bash
$ license-header-checker -v -a -r -i node_modules,docs ~/Dev/license_header.txt ~/Dev/projects/my-app js ts
```
