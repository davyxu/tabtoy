package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
)

func LoadAllTable() {
	var tab Table
	loadTableFromFile(&tab, "../all.pbb")

	fmt.Println(proto.MarshalTextString(&tab))
}

func loadTableFromFile(tab *Table, fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = proto.Unmarshal(data, tab)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func LoadSpecifiedTable() {
	var tabData Table
	loadTableFromFile(&tabData, "../ExampleData.pbb")

	fmt.Println("load specified table: ExampleData")
	for k, v := range tabData.ExampleData {
		fmt.Println(k, v)
	}

	var tabKV Table
	loadTableFromFile(&tabKV, "../ExampleKV.pbb")

	fmt.Println("load specified table: ExampleKV")
	for k, v := range tabKV.ExampleKV {
		fmt.Println(k, v)
	}

	// 将表格合并
	var total Table
	proto.Merge(&total, &tabData)
	proto.Merge(&total, &tabKV)
	fmt.Printf("merged table, data len: %d,  kv len: %d\n", len(total.ExampleData), len(total.ExampleKV))
}

func main() {

	LoadAllTable()

	LoadSpecifiedTable()
}
