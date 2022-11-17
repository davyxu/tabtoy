package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/pkg/errors"
	"strings"
)

func Compile(cp *model.Compiler) (ret error) {
	defer func() {

		switch err := recover().(type) {
		case *util.TableError:
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	if cp.DataFileGetter == nil {
		fileLoader := util.NewFileLoader(!cp.ParaLoading, cp.CacheDir)

		if cp.ParaLoading {
			for _, fileName := range cp.FileList {
				fileLoader.AddFile(fileName)
			}

			fileLoader.Commit()
		}
		cp.DataFileGetter = fileLoader
	}

	for _, fileName := range cp.FileList {
		file, err := cp.DataFileGetter.GetFile(fileName)
		if err != nil {
			return errors.Wrap(err, fileName)
		}

		loadDataTable(file, fileName, cp)
	}

	return
}

func loadDataTable(file util.TableFile, fileName string, cp *model.Compiler) {
	for _, sheet := range file.Sheets() {

		tab := model.NewDataTable()
		tab.HeaderType = sheet.Name()
		tab.FileName = fileName

		cp.Datas.AddDataTable(tab)

		maxCol := loadHeader(sheet, tab, cp.Types)

		// 遍历所有数据行
		for row := maxHeaderRow; ; row++ {

			if sheet.IsRowEmpty(row, maxCol+1) {
				break
			}

			// 读取每一行
			readOneRow(sheet, tab, row)
		}
	}
}

func readOneRow(sheet util.TableSheet, tab *model.DataTable, row int) bool {

	for _, header := range tab.Headers {

		if header.TypeInfo == nil {
			continue
		}

		// 浮点数用库取时，需要特殊处理
		isFloat := util.LanguagePrimitive(header.TypeInfo.FieldType, "go") == "float32"

		value := sheet.GetValue(row, header.Cell.Col, &util.ValueOption{ValueAsFloat: isFloat})

		// 首列带#时，本行忽略
		if header.Cell.Col == 0 && strings.HasPrefix(value, "#") {
			return false
		}

		cell := tab.MustGetCell(row-maxHeaderRow, header.Cell.Col)
		cell.Value = value
	}

	return true
}
