package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v4/compiler"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"os"
	"strings"
)

func parseInputFiles(g *model.Globals) {
	for _, name := range flag.Args() {
		raw := strings.Split(name, ":")
		if len(raw) == 2 {
			g.AddFile(raw[0], raw[1])
		} else {
			g.AddFile("", name)
		}
	}
}

func V4Entry() {
	report.Init()
	g := model.NewGlobals()
	g.ParaLoading = *paramPara

	parseInputFiles(g)

	if *paramUseCache {
		g.CacheDir = *paramCacheDir
		os.Mkdir(g.CacheDir, 0666)
	}

	err := compiler.Compile(g)
	if err != nil {
		goto Exit
	}

	return
Exit:
	report.Log.Errorln(err)
	os.Exit(1)
}
