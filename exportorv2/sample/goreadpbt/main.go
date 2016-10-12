package main

import (
	"fmt"
	"io/ioutil"

	"github.com/davyxu/tabtoy/exportorv2/sample/gamedef"
	"github.com/golang/protobuf/proto"
)

func main() {

	data, err := ioutil.ReadFile("../Config.pbt")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var config gamedef.Config

	err = proto.UnmarshalText(string(data), &config)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
