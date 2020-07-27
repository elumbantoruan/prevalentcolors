package datawriter

import (
	"errors"
	"os"
)

// DataWriter .
type DataWriter interface {
	Write(data []byte) error
}

// FileWriter is a concrete type of DataWriter which writes data into file
type FileWriter struct {
	Name string
	file *os.File
}

// NewFileWriter creates an instance of FileWriter
func NewFileWriter(name string) (*FileWriter, error) {
	if len(name) == 0 {
		return nil, errors.New("empty fileName")
	}
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		os.Create(name)
	}
	return &FileWriter{
		Name: name,
	}, nil
}

// Write writes data into a file
func (fw *FileWriter) Write(data []byte) error {

	file, err := os.OpenFile(fw.Name, os.O_APPEND|os.O_WRONLY, 06444)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	file.Close()

	return err
}
