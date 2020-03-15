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

func TestAddOperation(t *testing.T) {
	stats := NewStats()

	stats.AddOperation(&Operation{
		Action: LicenseAdded,
		Path:   "path1",
	})
	stats.AddOperation(&Operation{
		Action: LicenseReplaced,
		Path:   "path2",
	})
	stats.AddOperation(&Operation{
		Action: LicenseOk,
		Path:   "path3",
	})
	stats.AddOperation(&Operation{
		Action: LicenseOk,
		Path:   "path4",
	})

	assert.True(t, len(stats.Files[LicenseAdded]) == 1)
	assert.True(t, len(stats.Files[LicenseReplaced]) == 1)
	assert.True(t, len(stats.Files[LicenseOk]) == 2)
	assert.True(t, stats.Files[LicenseAdded][0] == "path1")
	assert.True(t, stats.Files[LicenseReplaced][0] == "path2")
	assert.True(t, stats.Files[LicenseOk][0] == "path3")
	assert.True(t, stats.Files[LicenseOk][1] == "path4")
}
