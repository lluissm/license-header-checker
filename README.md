![Header](images/header.png)

# license-header-checker

![Test Status](https://github.com/lluissm/license-header-checker/workflows/Test/badge.svg)

Multiplatform command line tool that **checks** whether the **license headers** are included in the **source files** of a project.

- It can **insert the license** to a file in case it does not exist (optional).
- It can **replace the license** of a file in case it is different than the target (optional).
- **Only the selected file extensions** will be **processed**.
- Specific **folders, files and paths can be ignored**.

_DISCLAIMER_

The tool looks for the keywords `license` or `copyright` inside the **first block comment of the file** to determine whether the file contains a valid license.

The block comment should be following the format `/* */`. Thus, while it does support the source files of languages like _Go, Rust, JavaScript, TypeScript, C, C++, Java, Swift, Kotlin and C#_, it does not support the file extensions that do not use this comment style.

Go build tags (or anything that is not a block comment that could be before the license) are respected when **replacing** the license.

## Command Usage

### Syntax

```bash
$ license-header-checker [-a] [-r] [-v] [-i path1,...] license-header-path src-path extensions...
```

### Options

```
  -a        Add the target license in case the file does not have any.
  -r        Replace the existing license by the target one in case they are different.
  -v        Be verbose during execution.
  -i        A comma separated list of the folders, files and/or paths that should be ignored.
            It does not support wildcards.
  -version  Display version number.
```

### Example

```bash
$ license-header-checker -v -a -r -i node_modules,client/assets ../license_header.txt . js ts
```

## Usage in CI

### GitHub Action example

```yml
name: License Check
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v3
      - name: Install license-header-checker
        run: curl -s https://raw.githubusercontent.com/lluissm/license-header-checker/master/install.sh | bash
      - name: Run license check
        run: ./bin/license-header-checker -a -r -i testdata ./license_header.txt . go && [[ -z `git status -s` ]]
```

## How to install

### Install script

For linux and MacOS systems you can use the install script.

To install the latest version:

```bash
$ curl -s https://raw.githubusercontent.com/lluissm/license-header-checker/master/install.sh | bash
```

To install a specific version (e.g., v.1.3.0):

```bash
$ curl -s https://raw.githubusercontent.com/lluissm/license-header-checker/master/install.sh | bash -s v1.3.0
```

### Binary packages

The binary packages for Linux, Windows and macOS are uploaded for each release and can be downloaded from the [releases](https://github.com/lluissm/license-header-checker/releases) page.

### Building from source

Provided you have Go and make installed, just type:

```bash
$ make install
```
