package v3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
)

func loadSymbols(globals *model.Globals, fileName string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	sheet := file.Sheets[0]

	for row := 1; ; row++ {

		if util.IsFullRowEmpty(sheet, row) {
			break
		}

		var objtype model.ObjectTypes
		objtype.Table = util.GetSheetValueString(sheet, row, 0)
		objtype.ObjectType = util.GetSheetValueString(sheet, row, 1)
		objtype.Name = util.GetSheetValueString(sheet, row, 2)
		objtype.FieldName = util.GetSheetValueString(sheet, row, 3)
		objtype.FieldType = util.GetSheetValueString(sheet, row, 4)
		objtype.DefaultValue = util.GetSheetValueString(sheet, row, 5)
		globals.Symbols.Types = append(globals.Symbols.Types, &objtype)

	}

	return nil
}
