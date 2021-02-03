package tabtoy

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Table interface {
	ResetData() error
	BuildData() error

	ResetTable(tableName string) error
	IndexTable(tableName string) error
}

func LoadFromFile(tab Table, filename string) error {

	// 根据需要从你的源数据读取，这里从指定文件名的文件读取
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadFromData(tab, data)
}

func LoadFromData(tab Table, data []byte) error {
	// 重置数据，这里会触发Prehandler
	err := tab.ResetData()

	if err != nil {
		return err
	}

	// 使用json反序列化
	err = json.Unmarshal(data, tab)

	if err != nil {
		return err
	}

	// 构建数据和索引，这里会触发PostHandler
	return tab.BuildData()
}

// 从单一表文件加载
func LoadTableFromFile(tab Table, filename string) error {

	// 根据需要从你的源数据读取，这里从指定文件名的文件读取
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	tableName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	return LoadTableromData(tab, tableName, data)
}

// 从单一表数据加载
func LoadTableromData(tab Table, tableName string, data []byte) error {
	err := tab.ResetTable(tableName)
	if err != nil {
		return err
	}

	// 使用json反序列化
	err = json.Unmarshal(data, tab)

	if err != nil {
		return err
	}

	return tab.IndexTable(tableName)
}
