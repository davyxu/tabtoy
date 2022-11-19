package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func WriteFile(filename string, data []byte) error {

	err := os.MkdirAll(filepath.Dir(filename), 0755)

	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(filename, data, 0666)
}

func ChangeExtension(filename, newExt string) string {

	file := filepath.Base(filename)

	return strings.TrimSuffix(file, path.Ext(file)) + newExt
}
