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
	ignore := []string{"node_modules", "test", "docs", "dont_like_this_file.py", "client/assets", "dashboard/public"}
	assert.True(t, ShouldIgnore("node_modules/index.js", ignore))
	assert.True(t, ShouldIgnore("test/my-test.cpp", ignore))
	assert.True(t, ShouldIgnore("dont_like_this_file.py", ignore))
	assert.True(t, ShouldIgnore("myapp/docs/index.html", ignore))
	assert.False(t, ShouldIgnore("node_modules/index.js", []string{}))
	assert.False(t, ShouldIgnore("node_modules/index.js", nil))
	assert.False(t, ShouldIgnore("src/testQ/my-file.cpp", ignore))
	assert.False(t, ShouldIgnore("src/TestData.java", ignore))
	assert.False(t, ShouldIgnore("src/test.py", ignore))
	assert.True(t, ShouldIgnore("dashboard/public/index.js", ignore))
	assert.True(t, ShouldIgnore("client/assets/index.js", ignore))
}
