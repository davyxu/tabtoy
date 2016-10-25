package main

import (
	"fmt"
	"io/ioutil"

	"github.com/davyxu/tabtoy/exportorv2/sample/gamedef"
	"github.com/davyxu/tabtoy/exportorv2/sample/table"
	"github.com/golang/protobuf/proto"
)

func main() {

	data, err := ioutil.ReadFile("../Config.pbt")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var config gamedef.Config

	// 加载配置
	err = proto.UnmarshalText(string(data), &config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 创建索引
	table.MakeIndex(&config)

	// 使用被索引的文件
	for index, v := range table.SampleByID {
		fmt.Println(index, v)
	}

}
