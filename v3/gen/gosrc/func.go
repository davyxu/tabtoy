package gosrc

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func KeyValueTypeNames(globals *model.Globals) (ret []string) {
	linq.From(globals.IndexList).WhereT(func(pragma *model.IndexDefine) bool {
		return pragma.Kind == model.TableKind_KeyValue
	}).SelectT(func(pragma *model.IndexDefine) string {

		return pragma.TableType
	}).Distinct().ToSlice(&ret)

	return
}

type TableIndices struct {
	Table     *model.DataTable
	FieldInfo *model.TypeDefine
}

func init() {
	UsefulFunc["GoTabTag"] = func(fieldType *model.TypeDefine) string {

		var sb strings.Builder

		var kv []string

		if fieldType.Name != "" {
			kv = append(kv, fmt.Sprintf("tb_name:\"%s\"", fieldType.Name))
		}

		if len(kv) > 0 {
			sb.WriteString("`")

			for _, s := range kv {
				sb.WriteString(s)
			}

			sb.WriteString("`")
		}

		return sb.String()
	}

	UsefulFunc["HasKeyValueTypes"] = func(globals *model.Globals) bool {
		return len(KeyValueTypeNames(globals)) > 0
	}

	UsefulFunc["GetKeyValueTypeNames"] = KeyValueTypeNames

	UsefulFunc["JsonTabOmit"] = func() string {
		return "`json:\"-\"`"
	}

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
