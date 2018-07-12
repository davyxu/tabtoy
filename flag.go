package main

import "flag"

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
	paramCombineStructName = flag.String("combinename", "Table", "combine struct name, code struct name")
	paramProtoOut          = flag.String("proto_out", "", "output protobuf define (*.proto)")
	paramPbtOut            = flag.String("pbt_out", "", "output proto text format (*.pbt)")
	paramLuaOut            = flag.String("lua_out", "", "output lua code (*.lua)")
	paramJsonOut           = flag.String("json_out", "", "output json format (*.json)")
	paramCSharpOut         = flag.String("csharp_out", "", "output c# class and deserialize code (*.cs)")
	paramGoOut             = flag.String("go_out", "", "output golang code (*.go)")
	paramBinaryOut         = flag.String("binary_out", "", "output binary format(*.bin)")
	paramTypeOut           = flag.String("type_out", "", "output table types(*.json)")
	paramCppOut            = flag.String("cpp_out", "", "output c++ format (*.cpp)")
)

// 特殊文件格式参数
var (
	paramProtoVersion = flag.Int("protover", 3, "output .proto file version, 2 or 3")

	paramLuaEnumIntValue = flag.Bool("luaenumintvalue", false, "use int type in lua enum value")
	paramLuaTabHeader    = flag.String("luatabheader", "", "output string to lua tab header")

	paramGenCSharpBinarySerializeCode = flag.Bool("cs_gensercode", true, "generate c# binary serialize code, default is true")
)
