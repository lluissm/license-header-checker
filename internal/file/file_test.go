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

package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileHasExtension(t *testing.T) {
	extensions := []string{".js", ".md"}
	assert.True(t, HasExtension("file.js", extensions))
	assert.True(t, HasExtension("readme.md", extensions))
	assert.False(t, HasExtension("index.html", extensions))
	assert.False(t, HasExtension("styles.css", extensions))
}

func TestFileShouldIgnoreFolder(t *testing.T) {
	assert.True(t, ShouldIgnore("node_modules/index.js", []string{"node_modules"}))
	assert.True(t, ShouldIgnore("src/test/fmt-test.cpp", []string{"test"}))
	assert.False(t, ShouldIgnore("src/mytest/my-file.cpp", []string{"test"}))
	assert.False(t, ShouldIgnore("src/testdata.py", []string{"test"}))
}

func TestFileShouldIgnoreFile(t *testing.T) {
	assert.True(t, ShouldIgnore("neural.py", []string{"neural.py"}))
	assert.True(t, ShouldIgnore("cmd/license-header-checker/main.go", []string{"main.go"}))
	assert.False(t, ShouldIgnore("neural.py", []string{"neural"}))
	assert.False(t, ShouldIgnore("cmd/license-header-checker/main.go", []string{"main"}))
}

func TestFileShouldIgnorePath(t *testing.T) {
	assert.True(t, ShouldIgnore("dashboard/public/index.js", []string{"dashboard/public"}))
	assert.True(t, ShouldIgnore("src/drivers/I2C.cpp", []string{"drivers/I2C.cpp"}))
	assert.False(t, ShouldIgnore("dashboard/publicity/index.js", []string{"dashboard/public"}))
	assert.False(t, ShouldIgnore("webclient/assets/index.js", []string{"client/assets"}))
	assert.False(t, ShouldIgnore("mysrc/drivers/I2C.cpp", []string{"src/drivers/I2C.cpp"}))
	assert.False(t, ShouldIgnore("src/drivers/I2C.cpp", []string{"/drivers/I2C"}))
}

func TestFileShouldIgnoreEmptyOrNil(t *testing.T) {
	assert.False(t, ShouldIgnore("src/utils/stringutils.c", []string{}))
	assert.False(t, ShouldIgnore("src/utils/stringutils.h", nil))
}

func TestFileShouldIgnoreUsesFullIgnoreList(t *testing.T) {
	assert.True(t, ShouldIgnore("node_modules/index.js", []string{"node_modules", "docs", "test"}))
	assert.True(t, ShouldIgnore("node_modules/index.js", []string{"docs", "node_modules", "test"}))
	assert.True(t, ShouldIgnore("node_modules/index.js", []string{"docs", "test", "node_modules"}))
}
