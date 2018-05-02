package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func LoadSymbols(globals *model.Globals, indexGetter FileGetter, fileName string) error {

	tabs, err := LoadTableData(indexGetter, fileName, "TableField")

	if err != nil {
		return err
	}

	var symbolTab model.SymbolTable
	for _, symbol := range table.CoreSymbols {
		symbolTab.AddField(symbol)
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TableField", &symbolTab)

		for row := 0; row < tab.RowCount(); row++ {

			var objtype table.TableField

			helper.ParseRow(&objtype, tab, row, &symbolTab)

			globals.Symbols.AddField(&objtype)
		}

	}

	return nil
}
