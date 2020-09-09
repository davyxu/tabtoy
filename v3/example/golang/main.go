package main

import (
	"fmt"
	tabtoy "github.com/davyxu/tabtoy/v3/api/golang"
	"os"
)

// 一次性加载所有表
func LoadAllTable() {

	var Tab = NewTable()

	// 表加载前清除之前的手动索引和表关联数据
	Tab.RegisterPreEntry(func(tab *Table) error {
		fmt.Println("tab pre load clear")
		return nil
	})

	// 表加载和构建索引后，需要手动处理数据的回调
	Tab.RegisterPostEntry(func(tab *Table) error {
		fmt.Println("tab post load done")
		fmt.Printf("%+v\n", tab.ExampleDataByID[200])

		fmt.Println("KV: ", tab.GetKeyValue_ExampleKV().ServerIP)
		return nil
	})

	err := tabtoy.LoadFromFile(Tab, "../json/table_gen.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
}

// 按指定表加载
func LoadSpecifiedTable() {
	var TabData = NewTable()
	err := tabtoy.LoadTableFromFile(TabData, "../jsondir/ExampleData.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("load specified table: ExampleData")
	for k, v := range TabData.ExampleDataByID {
		fmt.Println(k, v)
	}

	// 分表加载时, 不会触发pre/post Handler
	var TabKV = NewTable()
	err = tabtoy.LoadTableFromFile(TabKV, "../jsondir/ExampleKV.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("load specified table: ExampleKV")
	for k, v := range TabKV.ExampleKV {
		fmt.Println(k, v)
	}
}

func main() {
	LoadAllTable()

	LoadSpecifiedTable()
}
