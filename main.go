package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davyxu/golog"
	"github.com/davyxu/tabtoy/exportorv2"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/printer"
)

var log *golog.Logger = golog.New("main")

// 显示版本号
var paramVersion = flag.Bool("version", false, "Show version")

// 工作模式
var paramMode = flag.String("mode", "", "mode: exportorv2")

// 并发导出,提高导出速度, 输出日志会混乱
var paramPara = flag.Bool("para", false, "parallel export by your cpu count")

var paramProtoOut = flag.String("proto_out", "", "output protobuf define (*.proto)")
var paramPbtOut = flag.String("pbt_out", "", "output proto text format (*.pbt)")
var paramLuaOut = flag.String("lua_out", "", "output lua code (*.lua)")
var paramJsonOut = flag.String("json_out", "", "output json format (*.json)")
var paramCSharpOut = flag.String("csharp_out", "", "output c# class and deserialize code (*.cs)")
var paramGoOut = flag.String("go_out", "", "output golang code (*.go)")
var paramBinaryOut = flag.String("binary_out", "", "output binary format(*.bin)")
var paramTypeOut = flag.String("type_out", "", "output table types(*.json)")
var paramCombineStructName = flag.String("combinename", "", "combine struct name, code struct name")
var paramProtoVersion = flag.Int("protover", 3, "output .proto file version, 2 or 3")
var paramLanguage = flag.String("lan", "en_us", "set output language")

const Version = "2.8.3"

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Println(Version)
		return
	}

	switch *paramMode {
	case "exportorv2":

		g := printer.NewGlobals()

		if *paramLanguage != "" {
			if !i18n.SetLanguage(*paramLanguage) {
				log.Infof("language not support: %s", *paramLanguage)
			}
		}

		g.Version = Version

		for _, v := range flag.Args() {
			g.InputFileList = append(g.InputFileList, v)
		}

		g.ParaMode = *paramPara
		g.CombineStructName = *paramCombineStructName
		g.ProtoVersion = *paramProtoVersion

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

		if !exportorv2.Run(g) {
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
