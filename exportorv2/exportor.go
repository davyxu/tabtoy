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

		mergeFileList := strings.Split(inputFile, "+")

		for index, fileName := range mergeFileList {

			file := NewFile(fileName)

			if file == nil {
				return false
			}

			var mergeTarget string
			if len(mergeFileList) > 1 {
				mergeTarget = "--> " + filepath.Base(mergeFileList[0])
			}

			log.Infoln(filepath.Base(fileName), mergeTarget)

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

				// 只写入主文件的文件列表
				if file.Header != nil {
					fileObjList = append(fileObjList, file)
				}

				mainMergeFile = file
			} else {

				// 添加自文件
				mainMergeFile.mergeList = append(mainMergeFile.mergeList, file)

			}

		}

	}

	log.Infof("==========%s==========", i18n.String(i18n.Run_ExportSheetData))
	// 导出表格
	if !util.ParallelWorker(fileObjList, false, func(in interface{}) bool {

		file := in.(*File)

		log.Infoln(filepath.Base(file.FileName))

		dataModel := model.NewDataModel()

		tab := model.NewTable()
		tab.LocalFD = file.LocalFD

		// 主表
		if !file.ExportData(dataModel, nil) {
			return false
		}

		// 子表提供数据
		for _, mergeFile := range file.mergeList {

			log.Infoln(filepath.Base(mergeFile.FileName), "--->", filepath.Base(file.FileName))

			// 电子表格数据导出到Table对象
			if !mergeFile.ExportData(dataModel, file.Header) {
				return false
			}
		}

		if file.IsVertical() {
			mergeV(dataModel, tab)
		} else {
			mergeSTD(dataModel, tab)
		}

		// 整合类型信息和数据
		return g.AddContent(tab)

	}) {
		return false
	}

	// 根据各种导出类型, 调用各导出器导出
	return g.Print()
}
