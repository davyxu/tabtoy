package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 在包中定义外部可访问的表句柄
var Tab = NewTable()

// 重新加载指定文件名的表
func ReloadTable(filename string) {

	// 根据需要从你的源数据读取，这里从指定文件名的文件读取
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 重置数据，这里会触发Prehandler
	Tab.ResetData()

	// 使用json反序列化
	err = json.Unmarshal(data, Tab)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 构建数据和索引，这里会触发PostHandler
	Tab.BuildData()
}

func main() {

	// 表加载前清除之前的手动索引和表关联数据
	Tab.RegisterPreEntry(func(tab *Table) {

		fmt.Println("tab pre load clear")
	})

	// 表加载和构建索引后，需要手动处理数据的回调
	Tab.RegisterPostEntry(func(tab *Table) {

		fmt.Printf("%+v\n", tab.ExampleDataByID[200])

		fmt.Println("KV: ", tab.GetKeyValue_ExampleKV().ServerIP)
	})

	ReloadTable("../json/table_gen.json")

}
