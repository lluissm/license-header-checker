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

package process

import _ "embed"

//go:embed testdata/target_license_header.txt
var testTargetLicenseHeader string

//go:embed testdata/js_without_license.js
var testFileWithoutLicense string

//go:embed testdata/js_with_target_header.js
var testFileWithTargetLicense string

//go:embed testdata/js_with_different_header.js
var testFileWithDifferentLicense string

//go:embed testdata/go_with_target_header_and_build_tag.go
var testFileWithBuildTagsAndTargetLicense string

//go:embed testdata/go_with_different_header_and_build_tag.go
var testFileWithBuildTagsAndDifferentLicense string

//go:embed testdata/go_with_target_header_and_extra_comment.go
var testFileWithTargetLicenseAndExtraComments string

//go:embed testdata/py_target_license_header.txt
var testPythonTargetLicense string

//go:embed testdata/py_with_different_license.py
var testFileWithDifferentPythonTargetLicense string

//go:embed testdata/py_with_license.py
var testFileWithPythonTargetLicense string

//go:embed testdata/py_without_license.py
var testFileWithoutPythonTargetLicense string
