package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v3/compiler"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/binpak"
	"github.com/davyxu/tabtoy/v3/gen/cssrc"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/javasrc"
	"github.com/davyxu/tabtoy/v3/gen/jsontext"
	"github.com/davyxu/tabtoy/v3/gen/luasrc"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"os"
)

type V3GenEntry struct {
	name    string
	f       gen.GenFunc
	flagstr *string
}

// v3新增
var (
	paramIndexFile = flag.String("index", "", "input multi-files configs")

	paramUseGBKCSV = flag.Bool("use_gbkcsv", true, "use gbk format in csv file")
	paramMatchTag  = flag.String("matchtag", "", "match data table file tags in v3 Index file")

	v3GenList = []V3GenEntry{
		{"gosrc", gosrc.Generate, paramGoOut},
		{"jsontext", jsontext.Generate, paramJsonOut},
		{"luasrc", luasrc.Generate, paramLuaOut},
		{"cssrc", cssrc.Generate, paramCSharpOut},
		{"binpak", binpak.Generate, paramBinaryOut},
		{"javasrc", javasrc.Generate, paramJavaOut},
	}
)

func genFile(globals *model.Globals, entry V3GenEntry, c chan error) {
	filename := *entry.flagstr

	if data, err := entry.f(globals); err != nil {
		c <- err
	} else {

		report.Log.Infof("  [%s] %s", entry.name, filename)

		err = helper.WriteFile(filename, data)

		if err != nil {
			c <- err
		}
	}

	c <- nil
}

func GenFileByList(globals *model.Globals) error {

	var errList []chan error
	for _, entry := range v3GenList {

		if *entry.flagstr == "" {
			continue
		}

		c := make(chan error)
		errList = append(errList, c)
		go genFile(globals, entry, c)
	}

	for _, c := range errList {
		err := <-c
		if err != nil {
			return err
		}
	}

	return nil
}

func V3Entry() {
	globals := model.NewGlobals()
	globals.Version = Version
	globals.ParaLoading = *paramPara
	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName
	globals.GenBinary = *paramBinaryOut != ""
	globals.MatchTag = *paramMatchTag

	if globals.MatchTag != "" {
		report.Log.Infof("MatchTag: %s", globals.MatchTag)
	}

	idxloader := helper.NewFileLoader(true)
	idxloader.UseGBKCSV = *paramUseGBKCSV
	globals.IndexGetter = idxloader
	globals.UseGBKCSV = *paramUseGBKCSV

	var err error

	err = compiler.Compile(globals)

	if err != nil {
		goto Exit
	}

	report.Log.Debugln("Generate files...")
	err = GenFileByList(globals)
	if err != nil {
		goto Exit
	}

	return
Exit:
	report.Log.Errorln(err)
	os.Exit(1)
}
