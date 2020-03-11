/* MIT License

Copyright (c) 2020 Lluis Sanchez

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

import (
	"testing"

	"github.com/lsm-dev/license-header-checker/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func options() *Options {
	return &Options{
		Add:         true,
		Replace:     true,
		Path:        "path/to/src",
		Extensions:  []string{"cpp"},
		IgnorePaths: nil,
	}
}

func TestLicenseOk(t *testing.T) {
	filePath := "main.cpp"
	license := testdata.FakeTargetLicenseHeader
	fileContent := testdata.FakeFileWithTargetLicenseHeader
	options := options()
	handler := new(FakeHandlerSuccess)

	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseOk)
}

func TestAddLicense(t *testing.T) {
	filePath := "main.cpp"
	license := testdata.FakeTargetLicenseHeader
	fileContent := testdata.FakeFileWithoutLicense
	options := options()
	handler := new(FakeHandlerSuccess)
	handlerError := new(FakeHandlerError)

	options.Add = false
	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == SkippedAdd)

	options.Add = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseAdded)

	options.Add = true
	op = File(filePath, fileContent, license, options, handlerError)
	assert.True(t, op == OperationError)
}

func TestReplaceLicense(t *testing.T) {
	filePath := "main.cpp"
	license := testdata.FakeTargetLicenseHeader
	fileContent := testdata.FakeFileWithDifferentLicenseHeader
	options := options()
	handler := new(FakeHandlerSuccess)
	handlerError := new(FakeHandlerError)

	options.Replace = false
	op := File(filePath, fileContent, license, options, handler)
	assert.True(t, op == SkippedReplace)

	options.Replace = true
	op = File(filePath, fileContent, license, options, handler)
	assert.True(t, op == LicenseReplaced)

	options.Replace = true
	op = File(filePath, fileContent, license, options, handlerError)
	assert.True(t, op == OperationError)
}
