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

	var kvList, dataList model.DataTableList

	LoadIndex(globals, globals.IndexFile)

	for _, pragma := range globals.IndexList {

		fillTableType(pragma)

		switch pragma.TableMode {
		case table.TableMode_Data:

			tablist, err := LoadTableData(pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, tab.HeaderType, globals.Symbols)
				dataList.AddDataTable(tab)
			}

		case table.TableMode_Type:

			err = LoadSymbols(globals, pragma.TableFileName)

			if err != nil {
				return err
			}
		case table.TableMode_KeyValue:

			tablist, err := LoadTableData(pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				// 输入数据是按TableField格式写的，为了共享TableField字段
				ResolveHeaderFields(tab, "TableField", globals.Symbols)

				kvList.AddDataTable(tab)
			}
		}
	}

	// 合并所有的KV表行
	var mergedKV model.DataTableList
	mergeData(&kvList, &mergedKV, globals.Symbols)

	// 完整KV表转置为普通数据表
	for _, kvtab := range mergedKV.Data {
		ResolveHeaderFields(kvtab, kvtab.HeaderType, globals.Symbols)
		dataList.AddDataTable(transposeKVtoData(globals.Symbols, kvtab))
	}

	// 合并所有的数据表
	mergeData(&dataList, &globals.DataTableList, globals.Symbols)

	return nil
}
