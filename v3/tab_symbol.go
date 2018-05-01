package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func LoadSymbols(globals *model.Globals, fileName string) error {

	tabs, err := LoadTableData(fileName, "TableField")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		for row := 0; row < tab.RowCount(); row++ {

			var objtype table.TableField

			helper.ResolveRowByReflect(&objtype, tab, row)

			globals.Symbols.AddField(&objtype)
		}

	}

	return nil
}
