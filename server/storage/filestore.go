package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Storage interface that is implemented by storage providers
type Storage struct {
	filePath string
	data     map[string][]byte
}

// New creates a new file storage
func New(config ...Config) *Storage {
	cfg := configDefault(config...)
	data, err := loadFromFile(cfg.Folder)
	if err != nil {
		panic(err)
	}
	return &Storage{
		filePath: cfg.Folder,
		data:     data,
	}
}

// loadFromFile loads data from the specified file
func loadFromFile(filePath string) (map[string][]byte, error) {
	data := make(map[string][]byte)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return data, nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("error unmarshaling data: %w", err)
	}
	return data, nil
}

// Get value by key
func (s *Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val := s.data[key]
	return val, nil
}

// Set key with value
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	s.data[key] = val
	return s.saveToFile()
}

// Delete key by key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}

	delete(s.data, key)
	return s.saveToFile()
}

// Reset all keys
func (s *Storage) Reset() error {
	s.data = make(map[string][]byte)
	return s.saveToFile()
}

// saveToFile saves data to the file
func (s *Storage) saveToFile() error {
	bytes, err := json.Marshal(s.data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}
	return os.WriteFile(s.filePath, bytes, 0644)
}

// Close connection (no-op for file storage)
func (s *Storage) Close() error {
	return nil
}
