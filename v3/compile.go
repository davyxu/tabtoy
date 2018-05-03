package v3

import (
	"github.com/davyxu/tabtoy/v3/checker"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
)

type FileGetter interface {
	GetFile(filename string) (*xlsx.File, error)
}

func Compile(globals *model.Globals, indexGetter FileGetter) (ret error) {

	defer func() {

		switch err := recover().(type) {
		case *helper.ErrorObject:
			log.Errorf("%s", err.Error())
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	// TODO 更好的内建读取
	err := LoadTypeTable(globals, indexGetter, globals.BuiltinSymbolFile, true)

	if err != nil {
		return err
	}

	var kvList, dataList model.DataTableList

	LoadIndexTable(globals, indexGetter, globals.IndexFile)

	var loader FileGetter

	if globals.Para {
		// 缓冲文件
		asyncLoader := helper.NewAsyncFileLoader()

		for _, pragma := range globals.IndexList {
			asyncLoader.AddFile(pragma.TableFileName)
		}

		asyncLoader.Commit()

		loader = asyncLoader
	} else {
		loader = indexGetter
	}

	for _, pragma := range globals.IndexList {

		fillTableType(pragma)

		switch pragma.TableMode {
		case table.TableMode_Data:

			tablist, err := LoadDataTable(loader, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, tab.HeaderType, globals.Types)

				checker.CheckTypes(tab, globals.Types)

				dataList.AddDataTable(tab)
			}

		case table.TableMode_Type:

			err = LoadTypeTable(globals, loader, pragma.TableFileName, false)

			if err != nil {
				return err
			}
		case table.TableMode_KeyValue:

			tablist, err := LoadDataTable(loader, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, "TableKeyValue", globals.Types)

				checker.CheckTypes(tab, globals.Types)

				kvList.AddDataTable(tab)
			}
		}
	}

	// 合并所有的KV表行
	var mergedKV model.DataTableList
	mergeData(&kvList, &mergedKV, globals.Types)

	// 完整KV表转置为普通数据表
	for _, kvtab := range mergedKV.Data {
		ResolveHeaderFields(kvtab, kvtab.HeaderType, globals.Types)
		dataList.AddDataTable(transposeKVtoData(globals.Types, kvtab))
	}

	// 合并所有的数据表
	mergeData(&dataList, &globals.DataTableList, globals.Types)

	return nil
}
