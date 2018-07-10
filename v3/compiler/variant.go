package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
)

func loadVariantTables(globals *model.Globals, kvList, dataList *model.DataTableList) error {
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

				dataList.AddDataTable(tab)
			}

		case table.TableMode_Type:

			err := LoadTypeTable(globals.Types, globals.TableGetter, pragma.TableFileName, false)

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

	return nil
}
