package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
)

func Upgrade(globals *model.Globals) error {

	var err error

	globals.TargetTypesSheet = globals.AddTable("Type.xlsx")

	helper.WriteTypeTableHeader(globals.TargetTypesSheet)

	for _, inputFile := range globals.SourceFileList {

		for _, mergeFile := range strings.Split(inputFile, "+") {

			log.Infof("加载v2表 %s...", mergeFile)
			if err := loadTable(globals, mergeFile); err != nil {
				return err
			}
		}
	}

	log.Infoln("导出v3索引表 ...")
	err = ExportIndexTable(globals)
	if err != nil {
		return err
	}

	log.Infoln("导出v3类型表 ...")
	err = ExportTypes(globals)

	if err != nil {
		return err
	}

	return WriteOutput(globals)
}

func WriteOutput(globals *model.Globals) (ret error) {
	log.Debugln("输出v3表格:")

	globals.TargetTables.VisitAllTable(func(data *helper.MemFileData) bool {

		fullFileName := filepath.Join(globals.OutputDir, util.ChangeExtension(data.FileName, ".csv"))

		log.Infoln("\t", fullFileName)

		csvFile := helper.ConvertToCSV(data.File)

		csvFile.(*helper.CSVFile).Transform(helper.ConvUTF8ToGBK)

		ret = csvFile.Save(fullFileName)
		if ret != nil {
			return false
		}

		return true
	})

	return
}

func loadTypes(globals *model.Globals, sourceFile *xlsx.File, tabPragma *golexer.KVPair, fileName string) error {
	for _, sourceSheet := range sourceFile.Sheets {

		if sourceSheet.Name == "@Types" {

			if err := importTypes(globals, sourceSheet, tabPragma, fileName); err != nil {
				return err
			}
		}

	}

	return nil
}

func loadDatas(globals *model.Globals, sourceFile, targetFile *xlsx.File, tabPragma *golexer.KVPair, fileName string) error {
	for _, sourceSheet := range sourceFile.Sheets {

		if sourceSheet.Name != "@Types" {
			targetSheet, _ := targetFile.AddSheet(sourceSheet.Name)

			headerList := importDataHeader(globals, sourceSheet, targetSheet, tabPragma.GetString("TableName"))

			// 空表
			if len(headerList) == 0 {
				return nil
			}

			if err := importDatas(sourceSheet, targetSheet, headerList, fileName); err != nil {
				return err
			}

		}

	}

	return nil
}

func loadTable(globals *model.Globals, fileName string) error {

	sourceFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	tabPragma := golexer.NewKVPair()

	err = loadTypes(globals, sourceFile, tabPragma, fileName)

	if err != nil {
		return err
	}

	targetFile := xlsx.NewFile()

	err = loadDatas(globals, sourceFile, targetFile, tabPragma, fileName)

	if err != nil {
		return err
	}

	// 空表不输出
	if len(targetFile.Sheets) > 0 && len(targetFile.Sheets[0].Rows) > 0 {
		globals.AddTableByFile(fileName, tabPragma.GetString("TableName"), targetFile)
	}

	return nil
}
