package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v2"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/printer"
	"os"
)

func v2Entry() {
	g := printer.NewGlobals()

	if *paramLanguage != "" {
		if !i18n.SetLanguage(*paramLanguage) {
			log.Infof("language not support: %s", *paramLanguage)
		}
	}

	g.Version = Version_v2

	for _, v := range flag.Args() {
		g.InputFileList = append(g.InputFileList, v)
	}

	g.ParaMode = *paramPara
	g.CombineStructName = *paramCombineStructName
	g.ProtoVersion = *paramProtoVersion
	g.LuaEnumIntValue = *paramLuaEnumIntValue
	g.LuaTabHeader = *paramLuaTabHeader
	g.GenCSSerailizeCode = *paramGenCSharpBinarySerializeCode
	g.PackageName = *paramPackageName

	if *paramProtoOut != "" {
		g.AddOutputType("proto", *paramProtoOut)
	}

	if *paramPbtOut != "" {
		g.AddOutputType("pbt", *paramPbtOut)
	}

	if *paramJsonOut != "" {
		g.AddOutputType("json", *paramJsonOut)
	}

	if *paramLuaOut != "" {
		g.AddOutputType("lua", *paramLuaOut)
	}

	if *paramCSharpOut != "" {
		g.AddOutputType("cs", *paramCSharpOut)
	}

	if *paramGoOut != "" {
		g.AddOutputType("go", *paramGoOut)
	}

	if *paramCppOut != "" {
		g.AddOutputType("cpp", *paramCppOut)
	}

	if *paramBinaryOut != "" {
		g.AddOutputType("bin", *paramBinaryOut)
	}

	if *paramTypeOut != "" {
		g.AddOutputType("type", *paramTypeOut)
	}

	if !v2.Run(g) {
		os.Exit(1)
	}
}
