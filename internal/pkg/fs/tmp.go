package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrCreateFile = errors.New("failed to create file")
	ErrOpenFile   = errors.New("failed to open file")
	ErrReadFile   = errors.New("failed to read file")
	ErrWriteFile  = errors.New("failed to write to file")
	ErrDeleteFile = errors.New("failed to delete file")
	ErrCloseFile  = errors.New("failed to close file")
)

type FileInterface interface {
	Create(string, []byte) error
	Content(string) ([]byte, error)
	Delete(string) error
}

func New(dir string, prefix string) FileInterface {
	return &impl{dir: dir, prefix: prefix}
}

type impl struct {
	dir    string
	prefix string
}

func (f *impl) path(name string) string {
	return filepath.Join(f.dir, f.prefix+"-"+name)
}

func (f *impl) Create(name string, content []byte) error {
	file, err := os.Create(f.path(name))
	if err != nil {
		return ErrCreateFile
	}

	_, err = file.Write(content)
	if err != nil {
		return ErrWriteFile
	}

	if err := file.Close(); err != nil {
		return ErrCloseFile
	}

	return nil
}

func (f *impl) Content(name string) ([]byte, error) {
	file, err := os.Open(f.path(name))
	if err != nil {
		return nil, ErrOpenFile
	}

	readAll, err := io.ReadAll(file)
	if err != nil {
		return nil, ErrReadFile
	}

	if err := file.Close(); err != nil {
		return nil, ErrCloseFile
	}

	return readAll, nil
}

func (f *impl) Delete(name string) error {
	if err := os.Remove(f.path(name)); err != nil {
		return ErrDeleteFile
	}

	return nil
}
