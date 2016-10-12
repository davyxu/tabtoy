package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/davyxu/tabtoy/exportorv2/sample/gamedef"
)

func main() {

	data, err := ioutil.ReadFile("../Config.json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var config gamedef.Config

	err = json.Unmarshal(data, &config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
