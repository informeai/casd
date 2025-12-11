package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// UserHome returns the home directory of the current user.
func UserHome() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	if h := os.Getenv("USERPROFILE"); h != "" { // Windows
		return h
	}
	if h := os.Getenv("HOMEPATH"); h != "" { // Windows fallback
		drive := os.Getenv("HOMEDRIVE")
		return drive + h
	}
	return ""
}

// EnsureDir checks if a directory exists at the given path, and creates it if it does not exist.
func EnsureDir(path string) error {
	info, err := os.Stat(path)
	if os.IsExist(err) {
		return nil
	}
	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return os.ErrInvalid
	}

	return nil
}

// HashData returns the SHA-256 hash of the given data as a hexadecimal string.
func HashData(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// GetStoragePath returns the default path for the storage database.
func GetStoragePath() string {
	home := UserHome()
	return filepath.Join(home, ".casd", "storage.db")
}
