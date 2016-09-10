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
)

var log *golog.Logger = golog.New("main")

// 显示版本号
var paramVersion = flag.Bool("version", false, "Show version")

// 工作模式
var paramMode = flag.String("mode", "", "mode: exportorv1, exportorv2")

// 输入电子表格文件
var paramXlsFile = flag.String("xls", "XLS", "input excel file, use ',' splited file list by multipy files")

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
var paramProto3OutDir = flag.String("proto3_outdir", "", "[v2] output protobuf define v3 (*.proto)")
var paramProto2OutDir = flag.String("proto2_outdir", "", "[v2] output protobuf define v2 (*.proto)")
var paramPbtOutDir = flag.String("pbt_outdir", "", "[v2] output proto text format (*.pbt)")
var paramLuaOutDir = flag.String("lua_outdir", "", "[v2] output lua code (*.lua)")
var paramJsonOutDir = flag.String("json_outdir", "", "[v2] output json format (*.json)")
var paramCSharpOutDir = flag.String("csharp_outdir", "", "[v2] output c# class and deserialize code (*.cs)")
var paramBinaryOut = flag.String("binary_out", "", "[v2] input filename , output binary format(*.bin)")

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
			OutDir:        *paramOutDir,
		}) {
			goto Err
		}
	case "exportorv2":

		if !exportorv2.Run(exportorv2.Parameter{
			Version:       Version,
			InputFileList: flag.Args(),
			ParaMode:      *paramPara,
			Proto3OutDir:  *paramProto3OutDir,
			Proto2OutDir:  *paramProto2OutDir,
			PbtOutDir:     *paramPbtOutDir,
			JsonOutDir:    *paramJsonOutDir,
			LuaOutDir:     *paramLuaOutDir,
			CSharpOutDir:  *paramCSharpOutDir,
			BinaryFileOut: *paramBinaryOut,
		}) {
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
