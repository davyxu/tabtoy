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

	report.Log.Debugln("\n加载内建类型表:")
	err := LoadTypeTable(globals.Types, table.BuiltinTypes, "BuiltinTypes.xlsx", true)

	if err != nil {
		return err
	}

	report.Log.Debugln("\n加载索引表:")
	err = LoadIndexTable(globals, globals.IndexFile)

	if err != nil {
		return err
	}

	var kvList, dataList model.DataTableList

	// 加载多种表
	err = loadVariantTables(globals, &kvList, &dataList)

	if err != nil {
		return err
	}

	report.Log.Debugln("\n检查类型表:")
	CheckTypeTable(globals.Types)

	report.Log.Debugln("\n合并KV数据表:")

	// 合并所有的KV表行
	var mergedKV model.DataTableList
	mergeData(&kvList, &mergedKV, globals.Types)

	// 完整KV表转置为普通数据表
	for _, tab := range mergedKV.AllTables() {

		dataList.AddDataTable(transposeKVtoData(globals.Types, tab))
	}

	report.Log.Debugln("\n合并数据表:")

	// 合并所有的数据表
	mergeData(&dataList, &globals.Datas, globals.Types)

	report.Log.Debugln("\n完成:")

	return nil
}
