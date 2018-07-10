package v2tov3

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
)

func Upgrade(globals *model.Globals) error {

	var err error
	globals.TargetTypesSheet, err = globals.AddTable("Type.xlsx", "").AddSheet("Default")
	if err != nil {
		return err
	}

	helper.WriteTypeTableHeader(globals.TargetTypesSheet)

	for _, inputFile := range globals.SourceFileList {

		for _, mergeFile := range strings.Split(inputFile, "+") {

			if err := loadTable(globals, mergeFile); err != nil {
				return err
			}
		}
	}

	err = ExportIndexTable(globals)
	if err != nil {
		return err
	}

	err = ExportTypes(globals)

	if err != nil {
		return err
	}

	return WriteOutput(globals)
}

func markFileNameUpgrade(filename string) string {

	if filename != "Index.xlsx" && filename != "Type.xlsx" {
		return "Upgraded_"
	}

	return ""
}

func WriteOutput(globals *model.Globals) (ret error) {
	log.Debugln("输出v3表格:")

	globals.TargetTables.VisitAllTable(func(data *helper.MemFileData) bool {

		upgradeStr := markFileNameUpgrade(data.FileName)

		fullFileName := filepath.Join(globals.OutputDir, upgradeStr+data.FileName)

		log.Infoln("\t", fullFileName)

		ret = data.File.Save(fullFileName)
		if ret != nil {
			return false
		}

		return true
	})

	return
}

func loadTypes(globals *model.Globals, sourceFile *xlsx.File, tabPragma *golexer.KVPair) error {
	for _, sourceSheet := range sourceFile.Sheets {

		if sourceSheet.Name == "@Types" {

			if err := importTypes(globals, sourceSheet, tabPragma); err != nil {
				return err
			}
		}

	}

	return nil
}

func loadDatas(globals *model.Globals, sourceFile, targetFile *xlsx.File, tabPragma *golexer.KVPair) error {
	for _, sourceSheet := range sourceFile.Sheets {

		if sourceSheet.Name != "@Types" {
			targetSheet, _ := targetFile.AddSheet(sourceSheet.Name)

			headerList := importDataHeader(globals, sourceSheet, targetSheet, tabPragma.GetString("TableName"))

			// 空表
			if len(headerList) == 0 {
				return nil
			}

			if err := importDatas(globals, sourceSheet, targetSheet, headerList); err != nil {
				return err
			}

		}

	}

	return nil
}

func loadTable(globals *model.Globals, fileName string) error {
	sourceFile, err := globals.TableGetter.GetFile(fileName)
	if err != nil {
		return err
	}

	tabPragma := golexer.NewKVPair()

	loadTypes(globals, sourceFile, tabPragma)

	targetFile := xlsx.NewFile()

	loadDatas(globals, sourceFile, targetFile, tabPragma)

	// 空表不输出
	if len(targetFile.Sheets) > 0 && len(targetFile.Sheets[0].Rows) > 0 {
		globals.AddTableByFile(fileName, tabPragma.GetString("TableName"), targetFile)
	}

	return nil
}
