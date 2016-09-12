package exportorv2

import (
	"bytes"
	"path"

	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

type Parameter struct {
	Version           string
	InputFileList     []string
	ParaMode          bool
	PbtOutDir         string
	LuaOutDir         string
	JsonOutDir        string
	Proto3OutDir      string
	Proto2OutDir      string
	CSharpOutDir      string
	BinaryOutDir      string
	CombineStructName string // 不包含路径, 用作
}

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(param Parameter) bool {

	combineFile := printer.NewCombineFile(param.CombineStructName)

	if !util.ParallelWorker(param.InputFileList, param.ParaMode, func(inputFile string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(inputFile)
		if tab == nil {
			return false
		}

		if !combineFile.CombineType(inputFile, file.FileDescriptor) {
			return false
		}

		if !combineFile.WriteBinary(tab) {
			return false
		}

		/*
			if param.PbtOutDir != "" {

				filebase := util.ChangeExtension(inputFile, ".pbt")
				outputFile := path.Join(param.PbtOutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				if !printer.PrintPBT(tab, param.Version, outputFile) {
					return false
				}
			}

			if param.JsonOutDir != "" {

				filebase := util.ChangeExtension(inputFile, ".json")
				outputFile := path.Join(param.PbtOutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				if !printer.PrintJson(tab, param.Version, outputFile) {
					return false
				}
			}

			if param.LuaOutDir != "" {

				filebase := util.ChangeExtension(inputFile, ".lua")
				outputFile := path.Join(param.PbtOutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				if !printer.PrintLua(tab, param.Version, outputFile) {
					return false
				}
			}

			if param.CSharpOutDir != "" {

				filebase := util.ChangeExtension(inputFile, ".cs")

				outputFile := path.Join(param.CSharpOutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				bf := printer.PrintCSharp(file.TypeSet, nil, param.Version)
				if bf == nil {
					return false
				}

				if !bf.Write(outputFile) {
					return false
				}
			}

			if param.Proto3OutDir != "" {

				filebase := util.ChangeExtension(inputFile, ".proto")

				if file.TypeSet.Pragma.Proto3OutFileName != "" {
					filebase = file.TypeSet.Pragma.Proto3OutFileName
				}

				outputFile := path.Join(param.Proto3OutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				if !printer.PrintProto(file.TypeSet, 3, param.Version, outputFile) {
					return false
				}
			}

			if param.Proto2OutDir != "" {
				filebase := util.ChangeExtension(inputFile, ".proto")

				if file.TypeSet.Pragma.Proto2OutFileName != "" {
					filebase = file.TypeSet.Pragma.Proto2OutFileName
				}

				outputFile := path.Join(param.Proto2OutDir, filebase)

				log.Infof("%s%s\n", printIndent(2), filebase)

				if !printer.PrintProto(file.TypeSet, 2, param.Version, outputFile) {
					return false
				}
			}

		*/

		return true

	}) {
		return false
	}

	// 合并最终文件
	if param.BinaryOutDir != "" {

		outputFile := path.Join(param.BinaryOutDir, param.CombineStructName+".bin")

		// 输出合并后的C# XXFile结构
		if param.CSharpOutDir != "" {

			bf := combineFile.PrintCombineCSharp(param.Version)
			if bf == nil {
				return false
			}

			csharpOutputFile := path.Join(param.CSharpOutDir, param.CombineStructName+".cs")

			log.Infof("Combine C# Source: %s\n", csharpOutputFile)

			if !bf.Write(csharpOutputFile) {
				return false
			}

		}

		log.Infof("Combine Binary: %s\n", path.Join(param.CombineStructName+".bin"))

		if !combineFile.Write(outputFile) {
			return false
		}

	}

	return true
}
