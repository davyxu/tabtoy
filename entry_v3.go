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
	paramBuiltinSymbolFile = flag.Bool("builtinsymbol", false, "builtin symbols visible in table types")
	paramIndexFile         = flag.String("index", "", "input multi-files configs")

	v3GenList = []V3GenEntry{
		{gosrc.Generate, paramGoOut},
		{jsondata.Generate, paramJsonOut},
	}
)

func V3Entry() {
	globals := model.NewGlobals()
	globals.Version = Version_v3

	model.BuiltinSymbolsVisible = *paramBuiltinSymbolFile
	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName
	globals.Para = *paramPara

	globals.IndexGetter = new(helper.SyncFileLoader)

	if globals.Para {
		// 缓冲文件
		asyncLoader := helper.NewAsyncFileLoader()

		for _, pragma := range globals.IndexList {
			asyncLoader.AddFile(pragma.TableFileName)
		}

		asyncLoader.Commit()

		globals.TableGetter = asyncLoader
	} else {
		globals.TableGetter = globals.IndexGetter
	}

	err := v3.Compile(globals)

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
