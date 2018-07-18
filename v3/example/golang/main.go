package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {

	data, err := ioutil.ReadFile("../json_gen.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var tab Table
	err = json.Unmarshal(data, &tab)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", tab)

}
