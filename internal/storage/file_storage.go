package storage

import (
	"os"
	"path/filepath"
)

type FileStorage struct {
	baseDir string
}

func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{
		baseDir: baseDir,
	}
}

func (s *FileStorage) Put(id string, text []byte) error {
	filePath := filepath.Join(s.baseDir, id)
	return os.WriteFile(filePath, text, 0644)
}

func (s *FileStorage) Get(id string) ([]byte, error) {
	filePath := filepath.Join(s.baseDir, id)
	return os.ReadFile(filePath)
}

func (s *FileStorage) Del(id string) error {
	filePath := filepath.Join(s.baseDir, id)
	return os.Remove(filePath)
}

func (s *FileStorage) List() {}
