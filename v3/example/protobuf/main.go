package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
)

func main() {

	data, err := ioutil.ReadFile("./pb.bin")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var tab Table
	err = proto.Unmarshal(data, &tab)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(proto.MarshalTextString(&tab))
}
