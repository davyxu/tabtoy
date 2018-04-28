package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func Parse(globals *model.Globals) error {

	// TODO 更好的内建读取
	err := loadSymbols(globals, globals.BuiltinSymbolFile)

	if err != nil {
		return err
	}

	var kvlist model.DataTableList

	LoadIndex(globals, globals.IndexFile, func(pragma *table.TablePragma) error {

		switch pragma.TableType {
		case table.TableType_Data:

			tabName := getTableName(pragma)

			dataTable := globals.GetDataTable(tabName)

			if dataTable == nil {
				dataTable = model.NewDataTable()
				dataTable.Name = tabName
				globals.AddDataTable(dataTable)
			}

			err = LoadTableData(pragma.TableFileName, dataTable)

			if err != nil {
				return err
			}

		case table.TableType_Symbol:

			err = loadSymbols(globals, pragma.TableFileName)

			if err != nil {
				return err
			}
		case table.TableType_KeyValue:

			tabName := getTableName(pragma)

			kvtab := kvlist.GetDataTable(tabName)

			if kvtab == nil {
				kvtab = model.NewDataTable()
				kvtab.Name = tabName
				kvlist.AddDataTable(kvtab)
			}

			err = LoadTableData(pragma.TableFileName, kvtab)

			if err != nil {
				return err
			}

		}

		return nil
	})

	// kv转置
	for _, kvtab := range kvlist.Datas {
		ResolveHeaderFields(kvtab, "TableField", globals.Symbols)
		globals.AddDataTable(convertKVToData(globals.Symbols, kvtab))
	}

	for _, tab := range globals.Datas {
		ResolveHeaderFields(tab, tab.Name, globals.Symbols)
	}

	return nil
}
