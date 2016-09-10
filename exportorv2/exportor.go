package exportorv2

import (
	"bytes"
	"path"
	"path/filepath"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

type Parameter struct {
	Version       string
	InputFileList []string
	ParaMode      bool
	PbtOutDir     string
	LuaOutDir     string
	JsonOutDir    string
	Proto3OutDir  string
	Proto2OutDir  string
	CSharpOutDir  string
	BinaryFileOut string
}

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(param Parameter) bool {

	combineBinaryFile := printer.NewCombineBinaryFile()
	var combineFileTypes []*model.FieldDefine
	var combineFileNamespace string

	if !util.ParallelWorker(param.InputFileList, param.ParaMode, func(inputFile string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(inputFile)
		if tab == nil {
			return false
		}

		// 单个和合并二进制输出
		if param.BinaryFileOut != "" {

			rootName := file.TypeSet.Pragma.TableName

			bf := printer.PrintBinary(tab, rootName, param.Version)
			if bf == nil {
				return false
			}

			// 模块名字重复, 是无法输出的
			if !combineBinaryFile.Add(bf) {
				return false
			}

			// 有表格里描述的包名不一致, 无法合成最终的文件
			if combineFileNamespace != "" && combineFileNamespace != file.TypeSet.Pragma.Package {
				log.Errorf("combine file 'Package' in @Types diff: %s", inputFile)
				return false
			}

			combineFileNamespace = file.TypeSet.Pragma.Package

			combineFileTypes = append(combineFileTypes, file.TypeSet.FileType.Fields[0])

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

			bf := printer.PrintCSharp(file.TypeSet, param.Version)
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

		return true

	}) {
		return false
	}

	// 合并最终文件
	if param.BinaryFileOut != "" {

		filebase := filepath.Base(param.BinaryFileOut)

		combineName := util.ChangeExtension(param.BinaryFileOut, "File")

		// 输出合并后的C# XXFile结构
		if param.CSharpOutDir != "" {

			bf := printer.PrintCombineCSharp(combineFileTypes, param.Version, combineName, combineFileNamespace)
			if bf == nil {
				return false
			}

			csharpFileBase := util.ChangeExtension(param.BinaryFileOut, "File.cs")

			outputFile := path.Join(param.CSharpOutDir, csharpFileBase)

			log.Infof("Combine C# Source: %s\n", outputFile)

			if !bf.Write(outputFile) {
				return false
			}

		}

		log.Infof("Combine Binary: %s\n", filebase)

		if !combineBinaryFile.Write(param.BinaryFileOut) {
			return false
		}

	}

	return true
}
