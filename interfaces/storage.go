package interfaces

import "github.com/informeai/casd/dto"

// Storage defines the methods for storing and retrieving data chunks and formulas.
type Storage interface {
	PutChunk(hash string, data []byte) error
	GetChunk(hash string) ([]byte, error)
	DeleteChunk(hash string) error
	SaveFormula(key string, f *dto.Formula) error
	GetFormula(key string) (*dto.Formula, error)
	ListFormulas() (map[string]*dto.Formula, error)
	FilterFormulas(
		match func(f *dto.Formula) bool,
	) (map[string]*dto.Formula, error)
	Close() error
}
