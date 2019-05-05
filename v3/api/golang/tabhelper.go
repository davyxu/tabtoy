package tabtoy

import (
	"encoding/json"
	"io/ioutil"
)

type Table interface {
	ResetData() error
	BuildData() error
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
