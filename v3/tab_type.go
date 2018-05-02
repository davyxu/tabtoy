package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

var coreSymbols model.SymbolTable

func init() {

	for _, symbol := range table.CoreSymbols {
		coreSymbols.AddField(symbol)
	}
}

func LoadTypeTable(globals *model.Globals, indexGetter FileGetter, fileName string, builtin bool) error {

	tabs, err := LoadDataTable(indexGetter, fileName, "TableField")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TableField", &coreSymbols)

		for row := 0; row < tab.RowCount(); row++ {

			var objtype table.TableField

			helper.ParseRow(&objtype, tab, row, &coreSymbols)

			objtype.IsBuiltin = builtin

			globals.Symbols.AddField(&objtype)
		}

	}

	return nil
}
