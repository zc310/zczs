package pkg

import (
	"os"
	"path/filepath"
)

func ReadFile(name string) ([]byte, error) {
	b, err := os.ReadFile(filepath.Join(filepath.Dir(name), "历史", filepath.Base(name)))
	if err == nil {
		return b, nil
	}
	return os.ReadFile(name)
}
