package main

import (
	"fmt"

	"github.com/davyxu/tabtoy/v2/example/golang/table"
)

func main() {

	config := table.NewConfigTable()

	if err := config.Load("Config.json"); err != nil {
		panic(err)
	}

	for index, v := range config.SampleByID {
		fmt.Println(index, v)
	}

}
