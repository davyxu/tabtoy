package main

import (
	"flag"
	"github.com/davyxu/tabtoy/scanner"
	"strings"
)

///////////////////////////////////////////////
// mode: xls2pbt参数
///////////////////////////////////////////////
var paramSrcXls = flag.String("srcxls", "", "source xls file")

func runSyncHeaderMode() bool {

	headerFileName := *paramSrcXls

	// 打开提供文件头的XLS
	headerFile := scanner.NewFile(nil)

	if !headerFile.Open(headerFileName) {
		return false
	}

	// 遍历所有提供数据体的XLS
	for _, bodyFileName := range flag.Args() {

		// 打开数据体的XLS
		bodyfile := scanner.NewFile(nil)

		if !bodyfile.Open(bodyFileName) {
			return false
		}

		// 创建合并文件
		outFile := scanner.NewFile(nil)

		if !writeHeader(headerFile, outFile) {
			return false
		}

		if !writeBody(bodyfile, outFile) {
			return false
		}

		outFile.Raw.Save(bodyFileName + "_out")

	}

	return true

}

func getSheet(srcSheet *scanner.Sheet, out *scanner.File) *scanner.Sheet {
	outSheet, ok := out.SheetMap[srcSheet.Name]
	if !ok {
		rawSheet, _ := out.Raw.AddSheet(srcSheet.Name)

		outSheet = out.Add(rawSheet)
	}

	return outSheet

}

func writeHeader(src, out *scanner.File) bool {

	for _, srcSheet := range src.Sheets {

		if !srcSheet.ParseProtoField() {
			log.Errorf("src file '%s' lost proto header", src.FileName)
			return false
		}

		outSheet := getSheet(srcSheet, out)

		for cursor := 0; cursor < scanner.DataIndex_DataBegin; cursor++ {

			outRow := outSheet.AddRow()

			for index := 0; index < len(srcSheet.FieldHeader); index++ {

				outRow.AddCell().Value = srcSheet.GetCellData(cursor, index)
			}
		}
	}

	return true
}

func writeBody(src, out *scanner.File) bool {

	for _, srcSheet := range src.Sheets {

		// 保证有头, 虽然不读
		if !srcSheet.ParseProtoField() {
			log.Errorf("src file '%s' lost proto header", src.FileName)
			return false
		}

		outSheet := getSheet(srcSheet, out)

		for cursor := scanner.DataIndex_DataBegin; ; cursor++ {

			// 第一列是空的，结束
			if strings.TrimSpace(srcSheet.GetCellData(cursor, 0)) == "" {

				break
			}

			outRow := outSheet.AddRow()

			for index := 0; index < len(srcSheet.FieldHeader); index++ {

				outRow.AddCell().SetString(outRow.Cells[index].String())
			}
		}
	}

	return true
}
