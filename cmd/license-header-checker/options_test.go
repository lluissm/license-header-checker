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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorIfMissingArgs(t *testing.T) {
	args := []string{"license-header-checker"}
	_, err := parseOptions(args)
	assert.NotNil(t, err)

	args = []string{"license-header-checker", "license-path"}
	_, err = parseOptions(args)
	assert.NotNil(t, err)

	args = []string{"license-header-checker", "license-path", "source-path"}
	_, err = parseOptions(args)
	assert.NotNil(t, err)

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	_, err = parseOptions(args)
	assert.Nil(t, err)
}

func TestVersion(t *testing.T) {
	args := []string{"license-header-checker", "-version", "license-path", "source-path", "js"}
	options, _ := parseOptions(args)
	assert.True(t, options.ShowVersion)

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	options, _ = parseOptions(args)
	assert.False(t, options.ShowVersion)
}

func TestAdd(t *testing.T) {
	args := []string{"license-header-checker", "-a", "license-path", "source-path", "js"}
	options, _ := parseOptions(args)
	assert.True(t, options.Process.Add)

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	options, _ = parseOptions(args)
	assert.False(t, options.Process.Add)
}

func TestReplace(t *testing.T) {
	args := []string{"license-header-checker", "-r", "license-path", "source-path", "js"}
	options, _ := parseOptions(args)
	assert.True(t, options.Process.Replace)

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	options, _ = parseOptions(args)
	assert.False(t, options.Process.Replace)
}

func TestVerbose(t *testing.T) {
	args := []string{"license-header-checker", "-v", "license-path", "source-path", "js"}
	options, _ := parseOptions(args)
	assert.True(t, options.Verbose)

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	options, _ = parseOptions(args)
	assert.False(t, options.Verbose)
}

func TestPath(t *testing.T) {
	args := []string{"license-header-checker", "license-path", "source-path", "js", "ts"}
	options, _ := parseOptions(args)
	assert.True(t, options.Process.Path == "source-path")
}

func TestLicensePath(t *testing.T) {
	args := []string{"license-header-checker", "license-path", "source-path", "js", "ts"}
	options, _ := parseOptions(args)
	assert.True(t, options.Process.LicensePath == "license-path")
}

func TestExtensions(t *testing.T) {
	args := []string{"license-header-checker", "license-path", "source-path", "js", "ts"}
	options, _ := parseOptions(args)
	assert.True(t, len(options.Process.Extensions) == 2)
	assert.True(t, options.Process.Extensions[0] == ".js")
	assert.True(t, options.Process.Extensions[1] == ".ts")

	args = []string{"license-header-checker", "license-path", "source-path", "cpp"}
	options, _ = parseOptions(args)
	assert.True(t, len(options.Process.Extensions) == 1)
	assert.True(t, options.Process.Extensions[0] == ".cpp")
}

func TestIgnorePaths(t *testing.T) {
	args := []string{"license-header-checker", "-i", "node_modules,client/assets", "license-path", "source-path", "js"}
	options, _ := parseOptions(args)
	assert.True(t, len(options.Process.IgnorePaths) == 2)
	assert.True(t, options.Process.IgnorePaths[0] == "node_modules")
	assert.True(t, options.Process.IgnorePaths[1] == "client/assets")

	args = []string{"license-header-checker", "license-path", "source-path", "js"}
	options, _ = parseOptions(args)
	assert.True(t, len(options.Process.IgnorePaths) == 0)
}
