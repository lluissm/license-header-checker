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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HasExtension returns true if the file's extension is one of the provided ones
func HasExtension(path string, extensions []string) bool {
	fileExtension := filepath.Ext(path)
	for _, ext := range extensions {
		if fileExtension == ext {
			return true
		}
	}
	return false
}

// ShouldIgnore returns true if the path matches any of the ignore strings
func ShouldIgnore(path string, ignoreFolders []string) bool {
	segments := strings.Split(path, string(os.PathSeparator))
	for _, ignore := range ignoreFolders {
		if strings.Contains(ignore, string(os.PathSeparator)) {
			if strings.Contains(path, ignore) {
				return true
			}
			continue
		}
		for _, segment := range segments {
			if segment == ignore {
				return true
			}
		}
	}
	return false
}

// Replace remove the file and create a new one with the specified content
func Replace(path string, content string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed deleting the file: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed opening file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	if err = writer.Flush(); err != nil {
		return fmt.Errorf("failed writing to file: %w", err)
	}

	return nil
}
