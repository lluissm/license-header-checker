package process

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Handler handles adding and replacing license to a file
type handler struct{}

func (s *handler) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (s *handler) Walk(path string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(path, walkFn)
}

func (s *handler) ReplaceFileContent(filePath string, content string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed deleting the file: %w", err)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
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
