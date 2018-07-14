package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
)

func loadVariantTables(globals *model.Globals, kvList, dataList *model.DataTableList) error {
	report.Log.Debugln("Loading tables...")

	// 遍历索引里的每一行配置
	for _, pragma := range globals.IndexList {

		report.Log.Debugf("   (%s) %s", pragma.TableType, pragma.TableFileName)

		switch pragma.Kind {
		case model.TableKind_Data:
			tablist, err := LoadDataTable(globals.TableGetter, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {
				ResolveHeaderFields(tab, tab.HeaderType, globals.Types)

				CheckHeaderTypes(tab, globals.Types)

				dataList.AddDataTable(tab)
			}

		case model.TableKind_Type:

			err := LoadTypeTable(globals.Types, globals.TableGetter, pragma.TableFileName, false)

			if err != nil {
				return err
			}

		case model.TableKind_KeyValue:
			tablist, err := LoadDataTable(globals.TableGetter, pragma.TableFileName, pragma.TableType)

			if err != nil {
				return err
			}

			for _, tab := range tablist {

				ResolveHeaderFields(tab, "KVDefine", globals.Types)

				CheckHeaderTypes(tab, globals.Types)
				kvList.AddDataTable(tab)
			}

		}
	}

	return nil
}
