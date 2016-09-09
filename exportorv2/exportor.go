package exportorv2

import (
	"bytes"
	"path"
	"path/filepath"

	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

type Parameter struct {
	Version              string
	InputFileList        []string
	ParaMode             bool
	PbtOutDir            string
	LuaOutDir            string
	JsonOutDir           string
	Proto3OutDir         string
	Proto2OutDir         string
	CSharpOutDir         string
	BinaryOutDir         string
	CombineBinaryFileOut string
	CombineFileType      string
}

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(param Parameter) bool {

	var binaryByName = make(map[string][]byte)

	if !util.ParallelWorker(param.InputFileList, param.ParaMode, func(inputFile string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(inputFile)
		if tab == nil {
			return false
		}

		if param.BinaryOutDir != "" || param.CombineBinaryFileOut != "" {

			filebase := util.ChangeExtension(inputFile, ".bin")
			outputFile := path.Join(param.BinaryOutDir, filebase)

			log.Infof("%s%s\n", printIndent(2), filebase)

			rootName := file.TypeSet.Pragma.TableName

			fp := printer.PrintBinary(tab, rootName, param.Version)
			if fp == nil {
				return false
			}

			if param.CombineBinaryFileOut != "" {

				// 模块名字重复, 是无法输出的
				if _, ok := binaryByName[rootName]; ok {
					log.Errorln("duplicate table name in combine binary output:", rootName)
					return false
				}

				binaryByName[rootName] = fp.Data()

			} else {
				if !fp.Write(outputFile) {
					return false
				}
			}

		}

		if param.PbtOutDir != "" {

			filebase := util.ChangeExtension(inputFile, ".pbt")
			outputFile := path.Join(param.PbtOutDir, filebase)

			log.Infof("%s%s\n", printIndent(2), filebase)

			if !printer.PrintPBT(tab, file.TypeSet.Pragma.TableName, param.Version, outputFile) {
				return false
			}
		}

		if param.JsonOutDir != "" {

			filebase := util.ChangeExtension(inputFile, ".json")
			outputFile := path.Join(param.PbtOutDir, filebase)

			log.Infof("%s%s\n", printIndent(2), filebase)

			if !printer.PrintJson(tab, file.TypeSet.Pragma.TableName, param.Version, outputFile) {
				return false
			}
		}

		if param.LuaOutDir != "" {

			filebase := util.ChangeExtension(inputFile, ".lua")
			outputFile := path.Join(param.PbtOutDir, filebase)

			log.Infof("%s%s\n", printIndent(2), filebase)

			if !printer.PrintLua(tab, file.TypeSet.Pragma.TableName, param.Version, outputFile) {
				return false
			}
		}

		if param.CSharpOutDir != "" {

			filebase := util.ChangeExtension(inputFile, ".cs")

			outputFile := path.Join(param.CSharpOutDir, filebase)

			log.Infof("%s%s\n", printIndent(2), filebase)

			if !printer.PrintCSharp(file.TypeSet, param.Version, outputFile) {
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

		return true

	}) {
		return false
	}

	// 合并最终文件
	if param.CombineBinaryFileOut != "" {

		filebase := filepath.Base(param.CombineBinaryFileOut)

		log.Infof("%s%s\n", printIndent(2), filebase)

		//		if !printer.PrintBinary(tab, file.TypeSet.Pragma.TableName, param.Version, outputFile) {
		//			return false
		//		}
	}

	return true
}
