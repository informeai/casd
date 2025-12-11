package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/informeai/casd/dto"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketChunks   = "chunks"
	bucketFormulas = "formulas"
)

// DBBolt implements the interfaces.Storage using BoltDB as the underlying database.
type DBBolt struct {
	db *bolt.DB
}

// InitBolt initializes a BoltDB database at the given path with the specified file mode.
func InitBolt(path string, fileMode os.FileMode) (*DBBolt, error) {
	db, err := bolt.Open(path, fileMode, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	// cria buckets iniciais
	err = db.Update(func(tx *bolt.Tx) error {
		if _, e := tx.CreateBucketIfNotExists([]byte(bucketChunks)); e != nil {
			return e
		}
		if _, e := tx.CreateBucketIfNotExists([]byte(bucketFormulas)); e != nil {
			return e
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &DBBolt{db: db}, nil
}

// Close closes the BoltDB database.
func (b *DBBolt) Close() error {
	return b.db.Close()
}

// PutChunk salva bytes crus usando "hash" como chave
func (b *DBBolt) PutChunk(hash string, data []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketChunks))
		// evita sobrescrever (opcional): se quiser overwrite, remova o check
		if v := bk.Get([]byte(hash)); v != nil {
			return nil
		}
		return bk.Put([]byte(hash), data)
	})
}

// GetChunk recupera bytes crus pela chave "hash"
func (b *DBBolt) GetChunk(hash string) ([]byte, error) {
	var out []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketChunks))
		v := bk.Get([]byte(hash))
		if v == nil {
			return errors.New("chunk not found")
		}
		// copiar para não apontar para memória interna do mmap
		out = append([]byte(nil), v...)
		return nil
	})
	return out, err
}

// DeleteChunk removes a data chunk by its hash.
func (b *DBBolt) DeleteChunk(hash string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketChunks))
		return bk.Delete([]byte(hash))
	})
}

// SaveFormula saves a Formula struct using "key" as the identifier.
func (b *DBBolt) SaveFormula(key string, f *dto.Formula) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketFormulas))
		data, err := json.Marshal(f)
		if err != nil {
			return err
		}
		return bk.Put([]byte(key), data)
	})
}

// GetFormula retrieves a Formula struct by its identifier "key".
func (b *DBBolt) GetFormula(key string) (*dto.Formula, error) {
	var f dto.Formula
	err := b.db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketFormulas))
		v := bk.Get([]byte(key))
		if v == nil {
			return errors.New("formula not found")
		}
		return json.Unmarshal(v, &f)
	})
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FilterFormulas retrieves all Formula structs that match the provided filter function.
func (b *DBBolt) FilterFormulas(
	match func(f *dto.Formula) bool,
) (map[string]*dto.Formula, error) {

	results := map[string]*dto.Formula{}

	err := b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("formulas"))
		if b == nil {
			return fmt.Errorf("bucket formulas não encontrado")
		}

		return b.ForEach(func(k, v []byte) error {
			var f dto.Formula
			if err := json.Unmarshal(v, &f); err != nil {
				return err
			}

			if match(&f) {
				results[string(k)] = &f
			}

			return nil
		})
	})

	return results, err
}
