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
	"errors"
	"path/filepath"
)

// FakeIoHandlerSuccess will NOT return an error on its methods
type FakeIoHandlerSuccess struct{}

func (s *FakeIoHandlerSuccess) ReplaceFileContent(filePath string, fileContent string) error {
	return nil
}

func (s *FakeIoHandlerSuccess) ReadFile(filename string) ([]byte, error) {
	return nil, nil
}

func (s *FakeIoHandlerSuccess) Walk(path string, walkFn filepath.WalkFunc) error {
	return nil
}

// FakeIoHandlerError will return an error on its methods
type FakeIoHandlerError struct{}

func (s *FakeIoHandlerError) ReplaceFileContent(filePath string, fileContent string) error {
	return errors.New("error")
}

func (s *FakeIoHandlerError) ReadFile(filename string) ([]byte, error) {
	return nil, nil
}

func (s *FakeIoHandlerError) Walk(path string, walkFn filepath.WalkFunc) error {
	return nil
}
