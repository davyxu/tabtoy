package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/davyxu/golog"
	"github.com/davyxu/tabtoy/exportorv1"
	"github.com/davyxu/tabtoy/exportorv1/data"
	"github.com/davyxu/tabtoy/exportorv2"
	"github.com/davyxu/tabtoy/exportorv2/printer"
)

var log *golog.Logger = golog.New("main")

// 显示版本号
var paramVersion = flag.Bool("version", false, "Show version")

// 工作模式
var paramMode = flag.String("mode", "", "mode: exportorv1, exportorv2")

// 并发导出,提高导出速度, 输出日志会混乱
var paramPara = flag.Bool("para", false, "parallel export by your cpu count")

// ============================v1 版本参数============================

// 开启调试信息
var paramDebugLevel = flag.Int("debug", 0, "[v1] show debug info")

// 出现错误时暂停
var paramHaltOnError = flag.Bool("haltonerr", false, "[v1] halt on error")

// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
var paramPbFile = flag.String("pb", "PB", "[v1] input protobuf binary descript file, export by protoc-gen-meta plugins")

// 输出文件夹
var paramOutDir = flag.String("outdir", "OUT_DIR", "[v1] output directory")

// 补丁文件
var paramPatch = flag.String("patch", "", "[v1] patch input files then output")

// 输出文件格式
var paramFormat = flag.String("fmt", "pbt", "[v1] output file format, support 'pbt', 'lua' ")

// ============================v2 版本参数============================
var paramProtoOut = flag.String("proto_out", "", "[v2] output protobuf define (*.proto)")
var paramPbtOut = flag.String("pbt_out", "", "[v2] output proto text format (*.pbt)")
var paramLuaOut = flag.String("lua_out", "", "[v2] output lua code (*.lua)")
var paramJsonOut = flag.String("json_out", "", "[v2] output json format (*.json)")
var paramCSharpOut = flag.String("csharp_out", "", "[v2] output c# class and deserialize code (*.cs)")
var paramGoOut = flag.String("go_out", "", "[v2] output golang index code (*.go)")
var paramBinaryOut = flag.String("binary_out", "", "[v2] input filename , output binary format(*.bin)")
var paramCombineStructName = flag.String("combinename", "", "[v2] combine struct name, code struct name")
var paramProtoVersion = flag.Int("protover", 3, "[v2] output .proto file version, 2 or 3")

const Version = "2.0.0"

func main() {

	flag.Parse()

	// 版本
	if *paramVersion {
		fmt.Println(Version)
		return
	}

	switch *paramMode {
	case "xls2pbt", "exportorv1":
		// 调试信息挂接命令行
		data.DebuggingLevel = *paramDebugLevel
		if !exportor.Run(exportor.Parameter{
			InputFileList: flag.Args(),
			PBFile:        *paramPbFile,
			PatchFile:     *paramPatch,
			Format:        *paramFormat,
			ParaMode:      *paramPara,
		}) {
			goto Err
		}
	case "exportorv2":

		g := printer.NewGlobals()

		g.Version = Version
		g.InputFileList = flag.Args()
		g.ParaMode = *paramPara
		g.CombineStructName = *paramCombineStructName
		g.ProtoVersion = *paramProtoVersion

		if *paramProtoOut != "" {
			g.AddOutputType(".proto", *paramProtoOut)
		}

		if *paramPbtOut != "" {
			g.AddOutputType(".pbt", *paramPbtOut)
		}

		if *paramJsonOut != "" {
			g.AddOutputType(".json", *paramJsonOut)
		}

		if *paramLuaOut != "" {
			g.AddOutputType(".lua", *paramLuaOut)
		}

		if *paramCSharpOut != "" {
			g.AddOutputType(".cs", *paramCSharpOut)
		}

		if *paramGoOut != "" {
			g.AddOutputType(".go", *paramGoOut)
		}

		if *paramBinaryOut != "" {
			g.AddOutputType(".bin", *paramBinaryOut)
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

	if *paramHaltOnError {
		halt()
	}

	os.Exit(1)
	return

}

func halt() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadLine()
}
