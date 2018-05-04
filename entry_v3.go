package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v3"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/jsondata"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"os"
)

type V3GenEntry struct {
	f    gen.GenFunc
	name *string
}

// v3新增
var (
	paramBuiltinSymbolFile = flag.String("builtinsymbol", "", "input builtin symbol files describe types")
	paramIndexFile         = flag.String("index", "", "input multi-files configs")

	v3GenList = []V3GenEntry{
		{gosrc.Generate, paramGoOut},
		{jsondata.Generate, paramJsonOut},
	}
)

func V3Entry() {
	globals := model.NewGlobals()
	globals.Version = Version_v3
	globals.BuiltinSymbolFile = *paramBuiltinSymbolFile
	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName
	globals.Para = *paramPara

	// 内建build时，输出所有内置symbols
	if globals.BuiltinSymbolFile == "BuiltinTypes.xlsx" {
		model.UseAllBuiltinSymbols = true
	}

	err := v3.Compile(globals, new(helper.SyncFileLoader))

	if err != nil {
		report.Log.Errorln(err)
		os.Exit(1)
	}

	for _, entry := range v3GenList {

		if *entry.name == "" {
			continue
		}

		filename := *entry.name

		if data, err := entry.f(globals); err != nil {
			report.Log.Errorln(err)
			os.Exit(1)
		} else {

			report.Log.Infoln(filename)

			err = helper.WriteFile(filename, data)

			if err != nil {
				report.Log.Errorln(err)
				os.Exit(1)
			}

		}
	}
}
