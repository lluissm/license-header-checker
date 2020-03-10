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

package header

import (
	"strings"
	"testing"

	"github.com/lsm-dev/license-header-checker/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestContainsLicense(t *testing.T) {
	assert.True(t, ContainsLicense(testdata.FakeFileWithTargetLicenseHeader()))
	assert.True(t, ContainsLicense(testdata.FakeFileWithDifferentLicenseHeader()))
	assert.False(t, ContainsLicense(testdata.FakeFileWithoutLicense()))
}

func TestExtractHeader(t *testing.T) {
	expected := strings.TrimSpace(testdata.FakeTargetLicenseHeader())
	input := testdata.FakeFileWithTargetLicenseHeader()
	output := Extract(input)
	assert.True(t, output == expected)

	expected = "/* copyright */"
	input = "/* copyright */\nlorem ipsum dolor sit amet"
	output = Extract(input)
	assert.True(t, output == expected)
}

func TestRemoveHeader(t *testing.T) {
	expected := testdata.FakeFileWithoutLicense()

	input := testdata.FakeFileWithTargetLicenseHeader()
	output := Remove(input)
	assert.True(t, output == expected)

	input = testdata.FakeFileWithDifferentLicenseHeader()
	output = Remove(input)
	assert.True(t, output == expected)
}

func TestInsertHeader(t *testing.T) {
	expected := testdata.FakeFileWithTargetLicenseHeader()
	input := testdata.FakeFileWithoutLicense()
	header := testdata.FakeTargetLicenseHeader()
	output := Insert(input, header)
	assert.True(t, output == expected)
}

func TestReplaceHeader(t *testing.T) {
	expected := testdata.FakeFileWithTargetLicenseHeader()
	input := testdata.FakeFileWithDifferentLicenseHeader()
	header := testdata.FakeTargetLicenseHeader()
	output := Replace(input, header)
	assert.True(t, output == expected)
}
