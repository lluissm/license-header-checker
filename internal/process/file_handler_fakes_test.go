package process

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// FakeHandlerSuccess will NOT return an error when addLicense or replaceLicense are called
type FakeHandlerSuccess struct{}

func (s *FakeHandlerSuccess) ReplaceFileContent(filePath string, fileContent string) error {
	return nil
}

func (s *FakeHandlerSuccess) ReadFile(filename string) ([]byte, error) {
	return nil, nil
}

func (s *FakeHandlerSuccess) Walk(path string, walkFn filepath.WalkFunc) error {
	return nil
}

// FakeHandlerError will return an error when addLicense or replaceLicense are called
type FakeHandlerError struct{}

func (s *FakeHandlerError) ReplaceFileContent(filePath string, fileContent string) error {
	return errors.New("error")
}

func (s *FakeHandlerError) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (s *FakeHandlerError) Walk(path string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(path, walkFn)
}
