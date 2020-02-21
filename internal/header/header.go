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
)

func removeEmptyLinesAtStart(content string) string {
	for i, line := range strings.Split(strings.TrimSuffix(content, "\n"), "\n") {
		if len(line) > 0 {
			return content[i:]
		}
	}
	return ""
}

// ContainsLicense returns true if the content contains the words license or copyright in a header comment
func ContainsLicense(content string) bool {
	header := strings.ToLower(Extract(content))
	containsCopyright := strings.Contains(header, "copyright")
	containsLicense := strings.Contains(header, "license")
	return containsCopyright || containsLicense
}

// Extract returns the first block comment of the content (if any). Empty string otherwise.
func Extract(content string) (header string) {
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "*/") {
			header += line
			return
		}
		header += line + "\n"
	}
	return ""
}

// Remove removes the header from the content as well as potential empty lines between the header and body
func Remove(content string) string {
	header := Extract(content)
	body := strings.ReplaceAll(content, header, "")
	return removeEmptyLinesAtStart(body)
}

// Insert inserts the provided header at the beginning of the content
func Insert(content, header string) string {
	return header + "\n\n" + content
}

// Replace removes the current header and inserts the provided one
func Replace(content, header string) (res string) {
	res = Remove(content)
	res = Insert(res, header)
	return res
}
