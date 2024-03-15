package main

import (
	"github.com/davyxu/tabtoy/v4/compiler"
	"github.com/davyxu/tabtoy/v4/gen"
	"github.com/davyxu/tabtoy/v4/gen/bindata"
	"github.com/davyxu/tabtoy/v4/gen/cssrc"
	"github.com/davyxu/tabtoy/v4/gen/gosrc"
	"github.com/davyxu/tabtoy/v4/gen/jsondata"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"os"
)

type V4Generator struct {
	name    string
	genFunc gen.OutputFunc
	param   *string
}

var (
	v4Generator = []*V4Generator{
		{name: "gosrc", genFunc: gosrc.OutputFile, param: paramGoOut},
		{name: "cssrc", genFunc: cssrc.OutputFile, param: paramCSharpOut},
		{name: "jsondata", genFunc: jsondata.OutputFile, param: paramJsonOut},
		{name: "bindata", genFunc: bindata.OutputFile, param: paramBinaryOut},

		{name: "jsondir", genFunc: jsondata.OutputDir, param: paramJsonDir},
		{name: "bindir", genFunc: bindata.OutputDir, param: paramBinaryDir},
	}
)

func v4GenFile(globals *model.Globals, gen *V4Generator, c chan error) {
	err := gen.genFunc(globals, *gen.param)
	if err != nil {
		c <- err
	} else {
		report.Log.Infof("  [%s] %s", gen.name, *gen.param)
	}

	c <- nil
}

func V4BatchGenFile(globals *model.Globals) error {
	var errList []chan error
	for _, entry := range v4Generator {

		if *entry.param == "" {
			continue
		}

		c := make(chan error)
		errList = append(errList, c)
		go v4GenFile(globals, entry, c)
	}

	for _, c := range errList {
		err := <-c
		if err != nil {
			return err
		}
	}

	return nil
}

func V4Entry() {
	report.Init()
	g := model.NewGlobals()
	g.ParaLoading = *paramPara
	g.PackageName = *paramPackageName
	g.CombineStructName = *paramCombineStructName
	g.GenBinary = *paramBinaryOut != "" || *paramBinaryDir != ""

	err := compiler.ParseIndexFile(g, *paramIndexFile)
	if err != nil {
		goto Exit
	}

	if *paramUseCache {
		g.CacheDir = *paramCacheDir
		os.Mkdir(g.CacheDir, 0666)
	}

	err = compiler.Compile(g)
	if err != nil {
		goto Exit
	}

	if g.PackageName == "" {
		report.Log.Errorln("require package name, use --package")
	}

	if g.CombineStructName == "" {
		report.Log.Errorln("require combine struct name, use --combinename")
	}

	report.Log.Infof("Generate files...")
	err = V4BatchGenFile(g)
	if err != nil {
		goto Exit
	}

	return
Exit:
	report.Log.Errorln(err)
	os.Exit(1)
}
