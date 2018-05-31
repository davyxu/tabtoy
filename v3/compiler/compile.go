package compiler

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
)

func Compile(globals *model.Globals) (ret error) {

	defer func() {

		switch err := recover().(type) {
		case *report.TableError:
			fmt.Printf("%s\n", err.Error())
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	err := LoadTypeTable(globals.Types, table.BuiltinTypes, "BuiltinTypes.xlsx", true)

	if err != nil {
		return err
	}

	LoadIndexTable(globals, globals.IndexFile)

	var kvList, dataList model.DataTableList

	// 遍历索引里的每一行配置
	for _, pragma := range globals.IndexList {

		// 自动填充表项
		fillTableType(pragma)

		switch pragma.TableMode {
		case table.TableMode_Data:

			tablist, err := LoadDataTable(globals.TableGetter, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, tab.HeaderType, globals.Types)

				CheckHeaderTypes(tab, globals.Types)

				dataList.AddDataTable(tab)
			}

		case table.TableMode_Type:

			err = LoadTypeTable(globals.Types, globals.TableGetter, pragma.TableFileName, false)

			if err != nil {
				return err
			}
		case table.TableMode_KeyValue:

			tablist, err := LoadDataTable(globals.TableGetter, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, "TableKeyValue", globals.Types)

				CheckHeaderTypes(tab, globals.Types)

				kvList.AddDataTable(tab)
			}
		}
	}

	CheckTypeTable(globals.Types)

	// 合并所有的KV表行
	var mergedKV model.DataTableList
	mergeData(&kvList, &mergedKV, globals.Types)

	// 完整KV表转置为普通数据表
	for _, kvtab := range mergedKV.AllTables() {
		ResolveHeaderFields(kvtab, kvtab.HeaderType, globals.Types)
		dataList.AddDataTable(transposeKVtoData(globals.Types, kvtab))
	}

	// 合并所有的数据表
	mergeData(&dataList, &globals.Datas, globals.Types)

	return nil
}
