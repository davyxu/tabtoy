package main

import (
	"flag"
	"fmt"
	"github.com/davyxu/golog"
	"github.com/davyxu/tabtoy/v2"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/printer"
	"github.com/davyxu/tabtoy/v3"
	"github.com/davyxu/tabtoy/v3/genfile/gosrc"
	"github.com/davyxu/tabtoy/v3/genfile/json"
	"github.com/davyxu/tabtoy/v3/model"
	"os"
)

var log = golog.New("main")

// 标准参数
var (
	// 显示版本号
	paramVersion = flag.Bool("version", false, "Show version")

	// 工作模式
	paramMode = flag.String("mode", "", "v2")

	// 并发导出,提高导出速度, 输出日志会混乱
	paramPara = flag.Bool("para", false, "parallel export by your cpu count")

	paramLanguage = flag.String("lan", "en_us", "set output language")
)

// 文件类型导出
var (
	paramPackageName       = flag.String("package", "", "override the package name in table @Types")
	paramCombineStructName = flag.String("combinename", "Config", "combine struct name, code struct name")
	paramProtoOut          = flag.String("proto_out", "", "output protobuf define (*.proto)")
	paramPbtOut            = flag.String("pbt_out", "", "output proto text format (*.pbt)")
	paramLuaOut            = flag.String("lua_out", "", "output lua code (*.lua)")
	paramJsonOut           = flag.String("json_out", "", "output json format (*.json)")
	paramCSharpOut         = flag.String("csharp_out", "", "output c# class and deserialize code (*.cs)")
	paramGoOut             = flag.String("go_out", "", "output golang code (*.go)")
	paramBinaryOut         = flag.String("binary_out", "", "output binary format(*.bin)")
	paramTypeOut           = flag.String("type_out", "", "output table types(*.json)")
)

// 特殊文件格式参数
var (
	paramProtoVersion = flag.Int("protover", 3, "output .proto file version, 2 or 3")

	paramLuaEnumIntValue = flag.Bool("luaenumintvalue", false, "use int type in lua enum value")
	paramLuaTabHeader    = flag.String("luatabheader", "", "output string to lua tab header")

	paramGenCSharpBinarySerializeCode = flag.Bool("cs_gensercode", true, "generate c# binary serialize code, default is true")
)

// v3新增
var (
	paramSymbolFile = flag.String("symbol", "", "input symbol files describe types")
)

const (
	Version_v2 = "2.8.10"
	Version_v3 = "3.0.0"
)

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Printf("%s, %s", Version_v2, Version_v3)
		return
	}

	switch *paramMode {
	case "v3":

		var globals model.Globals
		globals.Version = Version_v3
		globals.SymbolFile = *paramSymbolFile
		globals.PackageName = *paramPackageName

		for _, v := range flag.Args() {
			globals.InputFileList = append(globals.InputFileList, v)
		}

		err := v3.Parse(&globals)
		if err != nil {
			fmt.Println(err)
		}

		if *paramJsonOut != "" {
			json.Generate(&globals, *paramJsonOut)
		}

		if *paramJsonOut != "" {
			gosrc.Generate(&globals, *paramGoOut)
		}

	case "exportorv2", "v2":

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

		if *paramBinaryOut != "" {
			g.AddOutputType("bin", *paramBinaryOut)
		}

		if *paramTypeOut != "" {
			g.AddOutputType("type", *paramTypeOut)
		}

		if !v2.Run(g) {
			goto Err
		}
	default:
		fmt.Println("--mode not specify")
		goto Err
	}

	return

Err:

	os.Exit(1)
	return

}
