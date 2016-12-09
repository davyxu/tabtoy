package table

// 本文件来自于: github.com/davyxu/tabtoy
// 请不要进行任何形式的修改
// 使用时, 请将将此文件放入你的工程的table包内
// tabtoy输出的go代码请与本文件放在同一个包内

import (
	"encoding/json"
	"io/ioutil"
)

var (
	indexEntryByName = make(map[string]func(interface{}))
)

func LoadTableFromFile(filename string, content interface{}) error {

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, content)

	if err != nil {
		return err
	}

	MakeIndex(content)

	return nil
}

func RegisterIndexEntry(name string, callback func(interface{})) {

	if _, ok := indexEntryByName[name]; ok {
		panic("duplicate table index entry")
	}

	indexEntryByName[name] = callback
}

func MakeIndex(content interface{}) {

	for _, v := range indexEntryByName {
		v(content)
	}

}
