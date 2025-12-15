package services

import (
	"path/filepath"

	"github.com/informeai/casd/dto"
	"github.com/informeai/casd/interfaces"
	"github.com/informeai/casd/repository"
	"github.com/informeai/casd/utils"
)

// Storage provides methods to interact with the storage database.
type Storage struct {
	db interfaces.Storage
}

// NewStorage initializes and returns a new Storage instance.
func NewStorage() (*Storage, error) {
	home := utils.UserHome()
	storagePath := utils.GetStoragePath()
	err := utils.EnsureDir(filepath.Join(home, ".casd"))
	if err != nil {
		return nil, err
	}

	db, err := repository.InitBolt(storagePath, 0755)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

// Close closes the storage database.
func (s *Storage) Close() error {
	return s.db.Close()
}

// PutChunk stores a data chunk with the given hash.
func (s *Storage) PutChunk(hash string, data []byte) error {
	return s.db.PutChunk(hash, data)
}

// GetChunk retrieves a data chunk by its hash.
func (s *Storage) GetChunk(hash string) ([]byte, error) {
	return s.db.GetChunk(hash)
}

// DeleteChunk removes a data chunk by its hash.
func (s *Storage) DeleteChunk(hash string) error {
	return s.db.DeleteChunk(hash)
}

// SaveFormula stores a formula with the given key.
func (s *Storage) SaveFormula(key string, f *dto.Formula) error {
	return s.db.SaveFormula(key, f)
}

// GetFormula retrieves a formula by its key.
func (s *Storage) GetFormula(key string) (*dto.Formula, error) {
	return s.db.GetFormula(key)
}

// ListFormulas retrieves all stored formulas.
func (s *Storage) ListFormulas() (map[string]*dto.Formula, error) {
	return s.db.ListFormulas()
}

// FilterFormulas retrieves formulas that match the given criteria.
func (s *Storage) FilterFormulas(
	match func(f *dto.Formula) bool,
) (map[string]*dto.Formula, error) {
	return s.db.FilterFormulas(match)
}
