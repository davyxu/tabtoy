package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteFile(filename string, data []byte) error {

	os.MkdirAll(filepath.Dir(filename), 0755)

	return ioutil.WriteFile(filename, data, 0666)
}
