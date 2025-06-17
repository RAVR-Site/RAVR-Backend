package storage

import (
	"mime/multipart"
)

// Storage interface for file storage operations
type Storage interface {
	// SaveFile сохраняет файл и возвращает URL для доступа к нему
	SaveFile(file *multipart.FileHeader, directory string) (string, error)

	// GetFileURL возвращает URL для доступа к файлу
	GetFileURL(filename string) string

	// DeleteFile удаляет файл из хранилища
	DeleteFile(filename string) error
}
