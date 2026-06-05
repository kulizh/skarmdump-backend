package storage

import (
	"os"
	"path/filepath"
)

type Local struct {
	base string
}

func NewLocal(base string) *Local {
	os.MkdirAll(base, os.ModePerm)
	return &Local{base: base}
}

func (l *Local) Save(hash string, data []byte) error {
	path := filepath.Join(l.base, hash+".png")
	if _, err := os.Stat(path); err == nil {
		return nil // dedup
	}
	return os.WriteFile(path, data, 0644)
}
