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

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// shouldIgnoreExtension returns false only if the file's extension is one of the provided ones
func shouldIgnoreExtension(path string, extensions []string) bool {
	fileExtension := filepath.Ext(path)
	for _, ext := range extensions {
		if fileExtension == ext {
			return false
		}
	}
	return true
}

// shouldIgnore returns true if the path matches any of the paths to ignore
func shouldIgnorePath(path string, ignorePaths []string) bool {
	pathSegments := strings.Split(path, string(os.PathSeparator))
	for _, ignorePath := range ignorePaths {
		ignorePathSegments := strings.Split(ignorePath, string(os.PathSeparator))
		size := len(ignorePathSegments)
		lastSegment := len(pathSegments) - size
		for i := 0; i <= lastSegment; i++ {
			if reflect.DeepEqual(pathSegments[i:i+size], ignorePathSegments) {
				return true
			}
		}
	}
	return false
}
