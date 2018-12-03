package helper

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteFile(filename string, data []byte) error {

	err := os.MkdirAll(filepath.Dir(filename), 0755)

	if err != nil && !os.IsExist(err) {
		return err
	}

	return ioutil.WriteFile(filename, data, 0666)
}
