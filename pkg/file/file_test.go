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
	assert.True(t, HasExtension("file.js", extensions), "Should return true for .js")
	assert.True(t, HasExtension("readme.md", extensions), "Should return true for .md")
	assert.False(t, HasExtension("index.html", extensions), "Should return false for .html")
	assert.False(t, HasExtension("styles.css", extensions), "Should return false for .css")
}

func TestFileShouldIgnore(t *testing.T) {
	ignore := []string{"node_modules", "docs"}
	assert.True(t, ShouldIgnore("node_modules", ignore), "Should return true for node_modules")
	assert.True(t, ShouldIgnore("docs", ignore), "Should return true for docs")
	assert.False(t, ShouldIgnore("package", ignore), "Should return false for package")
	assert.False(t, ShouldIgnore("src", ignore), "Should return false for src")
}
