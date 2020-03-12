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

	"github.com/stretchr/testify/assert"
)

func TestShouldIgnoreExtension(t *testing.T) {
	extensions := []string{".js", ".md"}
	assert.False(t, shouldIgnoreExtension("file.js", extensions))
	assert.False(t, shouldIgnoreExtension("readme.md", extensions))
	assert.True(t, shouldIgnoreExtension("index.html", extensions))
	assert.True(t, shouldIgnoreExtension("styles.css", extensions))
}

func TestShouldIgnoreFolder(t *testing.T) {
	assert.True(t, shouldIgnorePath("node_modules/index.js", []string{"node_modules"}))
	assert.True(t, shouldIgnorePath("src/test/fmt-test.cpp", []string{"test"}))
	assert.False(t, shouldIgnorePath("src/mytest/my-file.cpp", []string{"test"}))
	assert.False(t, shouldIgnorePath("src/testdata.py", []string{"test"}))
}

func TestShouldIgnoreFile(t *testing.T) {
	assert.True(t, shouldIgnorePath("neural.py", []string{"neural.py"}))
	assert.True(t, shouldIgnorePath("cmd/license-header-checker/main.go", []string{"main.go"}))
	assert.False(t, shouldIgnorePath("neural.py", []string{"neural"}))
	assert.False(t, shouldIgnorePath("cmd/license-header-checker/main.go", []string{"main"}))
}

func TestShouldIgnorePath(t *testing.T) {
	assert.True(t, shouldIgnorePath("dashboard/public/index.js", []string{"dashboard/public"}))
	assert.True(t, shouldIgnorePath("src/drivers/I2C.cpp", []string{"drivers/I2C.cpp"}))
	assert.False(t, shouldIgnorePath("dashboard/publicity/index.js", []string{"dashboard/public"}))
	assert.False(t, shouldIgnorePath("webclient/assets/index.js", []string{"client/assets"}))
	assert.False(t, shouldIgnorePath("mysrc/drivers/I2C.cpp", []string{"src/drivers/I2C.cpp"}))
	assert.False(t, shouldIgnorePath("src/drivers/I2C.cpp", []string{"/drivers/I2C"}))
}

func TestFileShouldIgnoreEmptyOrNil(t *testing.T) {
	assert.False(t, shouldIgnorePath("src/utils/stringutils.c", []string{}))
	assert.False(t, shouldIgnorePath("src/utils/stringutils.h", nil))
}

func TestFileShouldIgnoreUsesFullIgnoreList(t *testing.T) {
	assert.True(t, shouldIgnorePath("node_modules/index.js", []string{"node_modules", "docs", "test"}))
	assert.True(t, shouldIgnorePath("node_modules/index.js", []string{"docs", "node_modules", "test"}))
	assert.True(t, shouldIgnorePath("node_modules/index.js", []string{"docs", "test", "node_modules"}))
}
