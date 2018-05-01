package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func Parse(globals *model.Globals) error {

	defer func() {

		switch err := recover().(type) {
		case *helper.ErrorObject:
			log.Errorf("%s", err.Error())
		case nil:
		default:
			panic(err)
		}

	}()

	// TODO 更好的内建读取
	err := LoadSymbols(globals, globals.BuiltinSymbolFile)

	if err != nil {
		return err
	}

	var kvlist model.DataTableList

	LoadIndex(globals, globals.IndexFile)

	for _, pragma := range globals.IndexList {

		switch pragma.TableMode {
		case table.TableMode_Data:

			tabName := getTableName(pragma)

			dataTable := globals.GetDataTable(tabName)

			if dataTable == nil {
				dataTable = model.NewDataTable()
				dataTable.Name = tabName
				dataTable.FileName = pragma.TableFileName
				globals.AddDataTable(dataTable)
			}

			err = LoadTableData(pragma.TableFileName, dataTable)

			if err != nil {
				return err
			}

		case table.TableMode_Type:

			err = LoadSymbols(globals, pragma.TableFileName)

			if err != nil {
				return err
			}
		case table.TableMode_KeyValue:

			tabName := getTableName(pragma)

			kvtab := kvlist.GetDataTable(tabName)

			if kvtab == nil {
				kvtab = model.NewDataTable()
				kvtab.Name = tabName
				kvtab.FileName = pragma.TableFileName
				kvlist.AddDataTable(kvtab)
			}

			err = LoadTableData(pragma.TableFileName, kvtab)

			if err != nil {
				return err
			}
		}
	}

	// kv转置
	for _, kvtab := range kvlist.Data {
		ResolveHeaderFields(kvtab, "TableField", globals.Symbols)
		globals.AddDataTable(convertKVToData(globals.Symbols, kvtab))
	}

	for _, tab := range globals.Data {
		ResolveHeaderFields(tab, tab.Name, globals.Symbols)
	}

	return nil
}
