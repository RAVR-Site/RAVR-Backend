package local

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// Storage represents a local file storage implementation
type Storage struct {
	basePath string
	baseURL  string // базовый URL для доступа к файлам
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath, baseURL string) (*Storage, error) {
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", basePath, err)
	}
	return &Storage{basePath: basePath, baseURL: baseURL}, nil
}

// SaveFile сохраняет файл и возвращает URL для доступа к нему
func (s *Storage) SaveFile(file *multipart.FileHeader, directory string) (string, error) {
	// Создаем директорию, если она не существует
	dirPath := filepath.Join(s.basePath, directory)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Генерируем уникальное имя файла с временной меткой
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(dirPath, uniqueFilename)

	// Создаем целевой файл
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Копируем содержимое
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// Возвращаем относительный путь для сохранения в базе данных
	return filepath.Join(directory, uniqueFilename), nil
}

// GetFileURL возвращает URL для доступа к файлу
func (s *Storage) GetFileURL(filename string) string {
	return fmt.Sprintf("%s/%s", s.baseURL, filename)
}

// DeleteFile удаляет файл из хранилища
func (s *Storage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.basePath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // файл уже не существует, считаем операцию успешной
	}

	return os.Remove(filePath)
}
