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

func TestFileShouldIgnore(t *testing.T) {
	// Folder name matches => TRUE
	assert.True(t, ShouldIgnore("node_modules/index.js", []string{"node_modules"}))
	assert.True(t, ShouldIgnore("src/test/fmt-test.cpp", []string{"test"}))
	assert.True(t, ShouldIgnore("myapp/docs/index.html", []string{"docs"}))
	// File name matches => TRUE
	assert.True(t, ShouldIgnore("neural.py", []string{"neural.py"}))
	assert.True(t, ShouldIgnore("cmd/license-header-checker/main.go", []string{"main.go"}))
	// Empty and/or nil ignore folders => FALSE
	assert.False(t, ShouldIgnore("src/utils/stringutils.c", []string{}))
	assert.False(t, ShouldIgnore("src/utils/stringutils.h", nil))
	// File path contains folder name but it is not a folder => FALSE
	assert.False(t, ShouldIgnore("src/testQ/my-file.cpp", []string{"test"}))
	assert.False(t, ShouldIgnore("src/TestData.java", []string{"test"}))
	assert.False(t, ShouldIgnore("test.py", []string{"test"}))
	// File path matches ignore path => TRUE
	assert.True(t, ShouldIgnore("dashboard/public/index.js", []string{"dashboard/public"}))
	assert.True(t, ShouldIgnore("client/assets/index.js", []string{"client/assets"}))
	assert.True(t, ShouldIgnore("src/drivers/I2C.cpp", []string{"drivers/I2C.cpp"}))
	// File path contains ignore path but it is not a real path => FALSE
	assert.False(t, ShouldIgnore("dashboard/publicity/index.js", []string{"dashboard/public"}))
	assert.False(t, ShouldIgnore("webclient/assets/index.js", []string{"client/assets"}))
	assert.False(t, ShouldIgnore("mysrc/drivers/I2C.cpp", []string{"src/drivers/I2C.cpp"}))
	assert.False(t, ShouldIgnore("src/drivers/I2C.cpp", []string{"/drivers/I2C"}))
}
