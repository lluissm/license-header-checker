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
	"strings"
)

// containsLicenseHeader returns true if the content contains the words license or copyright in a header comment
func containsLicenseHeader(content string) bool {
	header := strings.ToLower(extractHeader(content))
	containsCopyright := strings.Contains(header, "copyright")
	containsLicense := strings.Contains(header, "license")
	return containsCopyright || containsLicense
}

// extractHeader returns the first block comment of the content (if any). Empty string otherwise.
func extractHeader(content string) (header string) {
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "*/") {
			header += line
			return
		}
		header += line + "\n"
	}
	return ""
}

// removeHeader removes the header from the content as well as potential empty lines between the header and the body
func removeHeader(content string) string {
	header := extractHeader(content)
	body := strings.ReplaceAll(content, header, "")
	return strings.TrimLeft(body, "\n")
}

// insertHeader inserts the provided header at the beginning of the content separated by one empty line
func insertHeader(content, header string) string {
	return strings.TrimSpace(header) + "\n\n" + strings.TrimLeft(content, "\n")
}

// replaceHeader removes the current header and inserts the provided one
func replaceHeader(content, header string) (res string) {
	res = removeHeader(content)
	res = insertHeader(res, header)
	return res
}
