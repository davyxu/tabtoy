package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/golog"
	"os"
)

var log = golog.New("main")

const (
	Version_v2 = "2.8.12"
)

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Printf("%s", Version_v2)
		return
	}

	switch *paramMode {
	case "v3":
		V3Entry()
	case "exportorv2", "v2":
		V2Entry()
	case "v2tov3":
		V2ToV3Entry()
	default:
		fmt.Println("--mode not specify")
		os.Exit(1)
	}

}
