package exportorv2

import (
	"path/filepath"
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/exportorv2/printer"
	"github.com/davyxu/tabtoy/util"
)

func Run(g *printer.Globals) bool {

	if !g.PreExport() {
		return false
	}

	fileObjList := make([]interface{}, 0)

	log.Infof("==========%s==========", i18n.String(i18n.Run_CollectTypeInfo))

	// 合并类型
	for _, in := range g.InputFileList {

		inputFile := in.(string)

		var mainMergeFile *File

		for index, fileName := range strings.Split(inputFile, "+") {

			file := NewFile(fileName)

			if file == nil {
				return false
			}

			log.Infoln(filepath.Base(fileName))

			file.GlobalFD = g.FileDescriptor

			// 电子表格数据导出到Table对象
			if !file.ExportLocalType() {
				return false
			}

			// 主文件才写入全局信息
			if index == 0 {

				// 整合类型信息和数据
				if !g.AddTypes(file.LocalFD) {
					return false
				}

				// 没有
				if file.Header != nil {
					fileObjList = append(fileObjList, file)
				}

				mainMergeFile = file
			} else {

				mainMergeFile.mergeNext = file

			}

		}

	}

	log.Infof("==========%s==========", i18n.String(i18n.Run_ExportSheetData))
	// 导出表格
	if !util.ParallelWorker(fileObjList, false, func(in interface{}) bool {

		file := in.(*File)

		var tab *model.Table

		for file != nil {

			log.Infoln(filepath.Base(file.FileName))

			// 电子表格数据导出到Table对象
			thisTab := file.ExportData()
			if thisTab == nil {
				return false
			}

			if tab == nil {
				tab = thisTab
			} else {

				// 合并表格数据到主tab
				tab.Recs = append(tab.Recs, thisTab.Recs...)
			}

			file = file.mergeNext
		}

		// 整合类型信息和数据
		return g.AddContent(tab)

	}) {
		return false
	}

	// 根据各种导出类型, 调用各导出器导出
	return g.Print()
}
