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

	globals.TargetTypesSheet = globals.TargetDatas.Create("Types.xlsx")
	helper.WriteTypeTableHeader(globals.TargetTypesSheet)

	for _, inputFile := range globals.SourceFileList {

		for _, mergeFile := range strings.Split(inputFile, "+") {

			if err := loadTable(globals, mergeFile); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadTypes(globals *model.Globals, sourceFile *xlsx.File, tabPragma *golexer.KVPair) error {
	for _, sourceSheet := range sourceFile.Sheets {

		if sourceSheet.Name == "@Types" {

			if err := procTypes(globals, sourceSheet, tabPragma); err != nil {
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

			headerList := procDataHeader(globals, sourceSheet, targetSheet, tabPragma.GetString("TableName"))

			if err := procDatas(globals, sourceSheet, targetSheet, headerList); err != nil {
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

	targetFile := xlsx.NewFile()

	fileName = filepath.Base(fileName)

	globals.TargetDatas.AddFile(fileName, targetFile)

	tabPragma := golexer.NewKVPair()

	loadTypes(globals, sourceFile, tabPragma)
	loadDatas(globals, sourceFile, targetFile, tabPragma)

	return nil
}
