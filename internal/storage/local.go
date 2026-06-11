package storage

import (
	"os"
	"path/filepath"
)

type Local struct {
	base string
}

func NewLocal(base string) *Local {
	for _, dir := range []string{base, filepath.Join(base, "og"), filepath.Join(base, "re")} {
		os.MkdirAll(dir, os.ModePerm)
	}
	return &Local{base: base}
}

func (l *Local) SavePNG(hash string, data []byte) error {
	path := filepath.Join(l.base, hash+".png")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, data, 0644)
}

func (l *Local) SaveOG(hash string, data []byte) error {
	path := filepath.Join(l.base, "og", hash+".png")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, data, 0644)
}

func (l *Local) SaveRE(hash string, data []byte) error {
	path := filepath.Join(l.base, "re", hash+".png")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	return os.WriteFile(path, data, 0644)
}

func (l *Local) Exists(hash string) bool {
	path := filepath.Join(l.base, hash+".png")
	_, err := os.Stat(path)
	return err == nil
}

func (l *Local) Read(hash string) ([]byte, error) {
	path := filepath.Join(l.base, hash+".png")
	return os.ReadFile(path)
}

func (l *Local) ReadOG(hash string) ([]byte, error) {
	path := filepath.Join(l.base, "og", hash+".png")
	return os.ReadFile(path)
}

func (l *Local) ReadRE(hash string) ([]byte, error) {
	path := filepath.Join(l.base, "re", hash+".png")
	return os.ReadFile(path)
}
