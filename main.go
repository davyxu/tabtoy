package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/golog"
	"os"
)

var log = golog.New("main")

const (
	Version_v2 = "2.9.1"
	Version_v3 = "3.0.0"
)

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Printf("%s, %s", Version_v2, Version_v3)
		return
	}

	switch *paramMode {
	case "v3":
		V3Entry()
	case "exportorv2", "v2":
		V2Entry()
	default:
		fmt.Println("--mode not specify")
		os.Exit(1)
	}

}
