package main

import (
	"fmt"

	"github.com/davyxu/tabtoy/exportorv2/sample/table"
)

func main() {

	var config table.Config

	if err := table.LoadTableFromFile("../Config.json", &config); err != nil {
		panic(err)
	}

	for index, v := range table.SampleByID {
		fmt.Println(index, v)
	}
}
