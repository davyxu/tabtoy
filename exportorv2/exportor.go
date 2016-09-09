package exportorv2

import (
	"bytes"
	"path"

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
}

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(param Parameter) bool {
	return util.ParallelWorker(param.InputFileList, param.ParaMode, func(inputFile string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(inputFile)
		if tab == nil {
			return false
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

	})

}
