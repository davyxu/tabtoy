package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/tabtoy/v3"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/json"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
	"os"
	"path/filepath"
)

type V3GenFunc func(globals *model.Globals) (data []byte, err error)

type V3GenEntry struct {
	f    V3GenFunc
	name *string
}

// v3新增
var (
	paramBuiltinSymbolFile = flag.String("builtinsymbol", "", "input builtin symbol files describe types")
	paramIndexFile         = flag.String("index", "", "input multi-files configs")

	v3GenList = []V3GenEntry{
		{gosrc.Generate, paramGoOut},
		{json.Generate, paramJsonOut},
	}
)

func v3Entry() {
	globals := model.NewGlobals()
	globals.Version = Version_v3
	globals.BuiltinSymbolFile = *paramBuiltinSymbolFile
	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName

	err := v3.Compile(globals, new(helper.SyncFileLoader))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, entry := range v3GenList {

		if *entry.name == "" {
			continue
		}

		filename := *entry.name

		if data, err := entry.f(globals); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			os.MkdirAll(filepath.Dir(filename), 0755)

			fmt.Println(filename)
			err = ioutil.WriteFile(filename, data, 0666)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}
	}
}
