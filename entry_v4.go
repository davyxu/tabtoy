package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v4/compiler"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"os"
	"strings"
)

func parseInputFiles(cp *model.Compiler) {
	for _, name := range flag.Args() {
		raw := strings.Split(name, ":")

		var fileMeta model.FileMeta

		if len(raw) == 2 {
			fileMeta.FileName = raw[1]
			fileMeta.Mode = raw[0]
		} else {
			fileMeta.FileName = name
		}

		cp.InputFiles = append(cp.InputFiles, fileMeta)
	}
}

func V4Entry() {
	report.Init()
	cp := model.NewCompiler()
	cp.ParaLoading = *paramPara

	parseInputFiles(cp)

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
