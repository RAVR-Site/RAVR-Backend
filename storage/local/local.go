package local

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, os.ModePerm)
	return &LocalStorage{basePath}
}

func (l *LocalStorage) Save(path, filename string) error {
	dst := filepath.Join(l.basePath, filename)
	// file is already saved to uploads; additional logic can go here
	fmt.Printf("Saved file to %s", dst)
	return nil
}
