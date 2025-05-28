package local

import (
	"fmt"
	"os"
	"path/filepath"
)

// Storage represents a local file storage implementation
type Storage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) (*Storage, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}
	return &Storage{basePath}, nil
}

// Save saves a file to the local storage
func (s *Storage) Save(_, filename string) error {
	dst := filepath.Join(s.basePath, filename)
	// file is already saved to uploads; additional logic can go here
	fmt.Printf("Saved file to %s", dst)
	return nil
}
