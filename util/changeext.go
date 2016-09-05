package util

import (
	"path"
	"path/filepath"
	"strings"
)

func ChangeExtension(filename, newExt string) string {

	file := filepath.Base(filename)

	return strings.TrimSuffix(file, path.Ext(file)) + newExt
}
