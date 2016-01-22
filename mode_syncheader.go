package main

import (
	"flag"
	"path/filepath"

	"github.com/davyxu/tabtoy/scanner"
)

///////////////////////////////////////////////
// mode: xls2pbt参数
///////////////////////////////////////////////
var paramSrcXls = flag.String("srcxls", "", "source xls file")

func runSyncHeaderMode() bool {

	srcFileName := *paramSrcXls

	// 打开源文件
	srcFile := scanner.NewFile(srcFileName, nil)
	if srcFile == nil {
		return false
	}

	for _, tgtFileName := range flag.Args() {

		tgtfile := scanner.NewFile(tgtFileName, nil)
		if tgtfile == nil {
			return false
		}

		if !syncFile(srcFile, tgtfile, srcFileName, tgtFileName) {
			return false
		}

		tgtfile.Save(tgtFileName + "2")

	}

	return true

}

func syncFile(src, tgt *scanner.File, srcFileName, tgtFileName string) bool {

	log.Infof("%s -> %s\n", filepath.Base(srcFileName), filepath.Base(tgtFileName))

	for _, sheet := range src.Sheets {

		if tgtSheet, ok := tgt.SheetMap[sheet.Name]; ok {

			if !syncSheet(sheet, tgtSheet, srcFileName, tgtFileName) {
				return false
			}

		} else {
			log.Errorf("target file %s lost sheet, name: '%s'\n", tgtFileName, sheet.Name)
			return false
		}

	}

	return true
}

func syncSheet(src, tgt *scanner.Sheet, srcFileName, tgtFileName string) bool {

	log.Infof("	%s", src.Name)

	if !src.ParseProtoField() {
		log.Errorf("src file '%s' lost proto header", srcFileName)
		return false
	}

	for cursor := 0; cursor < scanner.DataIndex_DataBegin; cursor++ {

		for index := 0; index < len(src.FieldHeader); index++ {
			v := src.GetCellData(cursor, index)
			tgt.SetCellData(cursor, index, v)
		}
	}

	return true
}
