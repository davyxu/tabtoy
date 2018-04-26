package v3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
	"reflect"
)

func loadSymbols(globals *model.Globals, fileName string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	sheet := file.Sheets[0]

	// 探测类型表头长度
	var maxCol int
	for {
		strValue := util.GetSheetValueString(sheet, 0, maxCol)

		// 空列，终止
		if strValue == "" {
			break
		}

		maxCol++
	}

	for row := 1; ; row++ {

		if util.IsFullRowEmpty(sheet, row) {
			break
		}

		var objtype table.TypeField

		vobjtype := reflect.ValueOf(&objtype).Elem()

		for col := 0; col < maxCol; col++ {

			strValue := util.GetSheetValueString(sheet, row, col)

			fieldType := vobjtype.Field(col)

			StringToValue(strValue, fieldType.Addr().Interface())

		}
		globals.Symbols.AddField(&objtype)

	}

	return nil
}
