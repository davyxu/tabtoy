package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v4/compiler"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"os"
)

func V4Entry() {
	report.Init()
	cp := model.NewCompiler()
	cp.ParaLoading = *paramPara
	cp.FileList = flag.Args()
	if *paramUseCache {
		cp.CacheDir = *paramCacheDir
		os.Mkdir(cp.CacheDir, 0666)
	}

	err := compiler.Compile(cp)
	if err != nil {
		goto Exit
	}

	return
Exit:
	report.Log.Errorln(err)
	os.Exit(1)
}
