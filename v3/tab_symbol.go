package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func LoadSymbols(globals *model.Globals, fileName string) error {

	var symbolTable = model.NewDataTable()
	symbolTable.FileName = fileName
	symbolTable.Name = "TableField"
	err := LoadTableData(fileName, symbolTable)

	if err != nil {
		return err
	}

	for row := 0; row < symbolTable.RowCount(); row++ {

		var objtype table.TableField

		helper.ResolveRowByReflect(&objtype, symbolTable, row)

		globals.Symbols.AddField(&objtype)
	}

	return nil
}
