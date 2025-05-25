package storage

// Storage interface for image saving
type Storage interface {
	Save(path, filename string) error
}
