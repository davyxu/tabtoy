package gen

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

type TableIndices struct {
	Table     *model.DataTable
	FieldInfo *model.TypeDefine
}

func KeyValueTypeNames(globals *model.Globals) (ret []string) {
	linq.From(globals.IndexList).WhereT(func(pragma *model.IndexDefine) bool {
		return pragma.Kind == model.TableKind_KeyValue
	}).SelectT(func(pragma *model.IndexDefine) string {

		return pragma.TableType
	}).Distinct().ToSlice(&ret)

	return
}

func init() {
	UsefulFunc["HasKeyValueTypes"] = func(globals *model.Globals) bool {
		return len(KeyValueTypeNames(globals)) > 0
	}

	UsefulFunc["GetKeyValueTypeNames"] = KeyValueTypeNames

	UsefulFunc["GetIndices"] = func(globals *model.Globals) (ret []TableIndices) {

		for _, tab := range globals.Datas.AllTables() {

			// 遍历输入数据的每一列
			for _, header := range tab.Headers {

				// 输入的列头
				if header.TypeInfo == nil {
					continue
				}

				if header.TypeInfo.MakeIndex {

					ret = append(ret, TableIndices{
						Table:     tab,
						FieldInfo: header.TypeInfo,
					})
				}
			}
		}

		return

	}
}
