package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
)

func Compile(globals *model.Globals) (ret error) {

	defer func() {

		switch err := recover().(type) {
		case *report.TableError:
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	//report.Log.Debugln("\n加载内建类型表:")
	err := LoadTypeTable(globals.Types, table.BuiltinTypes, "BuiltinTypes.xlsx", true)

	if err != nil {
		return err
	}

	//report.Log.Debugln("\n加载索引表:")
	LoadIndexTable(globals, globals.IndexFile)

	var kvList, dataList model.DataTableList

	report.Log.Debugln("\n加载表:")

	// 遍历索引里的每一行配置
	for _, pragma := range globals.IndexList {

		switch pragma.TableMode {
		case table.TableMode_Data:
			tablist, err := LoadDataTable(globals.TableGetter, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, tab.HeaderType, globals.Types)

				CheckHeaderTypes(tab, globals.Types)

				report.Log.Debugln(tab.String())
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

				report.Log.Debugln(tab.String())
				kvList.AddDataTable(tab)
			}

		}
	}

	CheckTypeTable(globals.Types)

	report.Log.Debugln("\n合并KV数据表:")

	// 合并所有的KV表行
	var mergedKV model.DataTableList
	mergeData(&kvList, &mergedKV, globals.Types)

	// 完整KV表转置为普通数据表
	for _, tab := range mergedKV.AllTables() {

		dataList.AddDataTable(transposeKVtoData(globals.Types, tab))
	}

	report.Log.Debugln("\n合并所有数据表:")

	// 合并所有的数据表
	mergeData(&dataList, &globals.Datas, globals.Types)

	report.Log.Debugln("\n完成:")
	for _, tab := range globals.Datas.AllTables() {
		report.Log.Debugln(tab.String())
	}

	return nil
}
