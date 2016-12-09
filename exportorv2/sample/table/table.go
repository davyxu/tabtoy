package table

import (
	"encoding/json"
	"io/ioutil"
)

// 请将将此文件及所在包引入工程(不包含table_gen.go)

// table的索引入口函数

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
