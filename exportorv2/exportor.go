package exportorv2

import (
	"bytes"

	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(g *printer.Globals) bool {

	if !util.ParallelWorker(g.InputFileList, g.ParaMode, func(inputFile string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(inputFile)
		if tab == nil {
			return false
		}

		if !g.CombineType(inputFile, file.FileDescriptor) {
			return false
		}

		if !g.CombineData(tab) {
			return false
		}

		return true

	}) {
		return false
	}

	return g.Run()
}
