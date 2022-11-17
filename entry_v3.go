package main

import (
	"flag"
	"github.com/davyxu/tabtoy/build"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/compiler"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/bindata"
	"github.com/davyxu/tabtoy/v3/gen/cssrc"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/javasrc"
	"github.com/davyxu/tabtoy/v3/gen/jsondata"
	"github.com/davyxu/tabtoy/v3/gen/jsontype"
	"github.com/davyxu/tabtoy/v3/gen/luasrc"
	"github.com/davyxu/tabtoy/v3/gen/pbdata"
	"github.com/davyxu/tabtoy/v3/gen/pbsrc"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"os"
)

type V3GenEntry struct {
	name          string
	genSingleFile gen.GenSingleFile
	genCustom     gen.GenCustom
	param         *string
}

// v3新增
var (
	paramIndexFile = flag.String("index", "", "input multi-files configs")
	paramTagAction = flag.String("tag_action", "", "do action by tag selected target, format: action1:tag1+tag2|action2:tag1+tag3")

	v3GenList = []V3GenEntry{
		{name: "gosrc", genSingleFile: gosrc.Generate, param: paramGoOut},
		{name: "jsondata", genSingleFile: jsondata.Generate, param: paramJsonOut},
		{name: "jsontype", genSingleFile: jsontype.Generate, param: paramJsonTypeOut},
		{name: "luasrc", genSingleFile: luasrc.Generate, param: paramLuaOut},
		{name: "cssrc", genSingleFile: cssrc.Generate, param: paramCSharpOut},
		{name: "bindata", genSingleFile: bindata.Generate, param: paramBinaryOut},
		{name: "javasrc", genSingleFile: javasrc.Generate, param: paramJavaOut},
		{name: "pbsrc", genSingleFile: pbsrc.Generate, param: paramProtoOut},
		{name: "pbdata", genSingleFile: pbdata.Generate, param: paramPbBinaryOut},

		{name: "jsondir", genCustom: jsondata.Output, param: paramJsonDir},
		{name: "luadir", genCustom: luasrc.Output, param: paramLuaDir},
		{name: "binarydir", genCustom: bindata.Output, param: paramBinaryDir},
		{name: "pbdatadir", genCustom: pbdata.Output, param: paramPbBinaryDir},
	}
)

func genFile(globals *model.Globals, entry V3GenEntry, c chan error) {
	filename := *entry.param

	if entry.genSingleFile != nil {
		if data, err := entry.genSingleFile(globals); err != nil {
			c <- err
		} else {

			report.Log.Infof("  [%s] %s", entry.name, filename)

			err = util.WriteFile(filename, data)

			if err != nil {
				c <- err
			}
		}
	}

	if entry.genCustom != nil {
		if err := entry.genCustom(globals, *entry.param); err != nil {
			c <- err
		} else {
			report.Log.Infof("  [%s] %s", entry.name, filename)
		}
	}

	c <- nil
}

func GenFileByList(globals *model.Globals) error {

	var errList []chan error
	for _, entry := range v3GenList {

		if *entry.param == "" {
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
	report.Init()
	globals := model.NewGlobals()
	globals.Version = build.Version
	globals.ParaLoading = *paramPara
	if *paramUseCache {
		globals.CacheDir = *paramCacheDir
		os.Mkdir(globals.CacheDir, 0666)
	}
	globals.IndexFile = *paramIndexFile
	globals.PackageName = *paramPackageName
	globals.CombineStructName = *paramCombineStructName
	globals.GenBinary = *paramBinaryOut != "" || *paramBinaryDir != ""

	idxloader := util.NewFileLoader(true, globals.CacheDir)
	globals.IndexGetter = idxloader

	var err error
	if *paramTagAction != "" {
		globals.TagActions, err = model.ParseTagAction(*paramTagAction)
		if err != nil {
			goto Exit
		}
	}

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
