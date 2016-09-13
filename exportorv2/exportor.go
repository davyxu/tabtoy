package exportorv2

import (
	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

func Run(g *printer.Globals) bool {

	if !g.PreExport() {
		return false
	}

	if !util.ParallelWorker(g.InputFileList, g.ParaMode, func(inputFile string) bool {

		file := NewFile()

		// 电子表格数据导出到Table对象
		tab := file.Export(inputFile)
		if tab == nil {
			return false
		}

		// 整合类型信息和数据
		return g.AddContent(tab)

	}) {
		return false
	}

	// 根据各种导出类型, 调用各导出器导出
	return g.Print()
}
