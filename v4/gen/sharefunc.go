package gen

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

type TableIndices struct {
	Table     *model.DataTable
	FieldInfo *model.DataField
}

func KeyValueTypeNames(globals *model.Globals) (ret []string) {

	for _, tab := range globals.Datas.AllTables() {
		if tab.Mode == "KV" {
			ret = append(ret, tab.HeaderType)
		}
	}

	return
}

func WrapSingleValue(globals *model.Globals, valueType *model.DataField, value string) string {
	switch {
	case valueType.FieldType == "string": // 字符串
		return util.StringWrap(util.StringEscape(value))
	case valueType.FieldType == "float":

		if value == "" {
			return util.FetchDefaultValue(valueType.FieldType)
		}

		return value
	case globals.Types.IsEnumKind(valueType.FieldType): // 枚举
		return globals.Types.ResolveEnumValue(valueType.FieldType, value)
	case valueType.FieldType == "bool":

		v, _ := util.ParseBool(value)
		if v {
			return "true"
		}

		return "false"
	}

	if value == "" {
		return util.FetchDefaultValue(valueType.FieldType)
	}

	return value
}

func GetIndicesByTable(tab *model.DataTable) (ret []TableIndices) {
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

	return
}

func GetIndices(globals *model.Globals) (ret []TableIndices) {
	for _, tab := range globals.Datas.AllTables() {
		ret = append(ret, GetIndicesByTable(tab)...)
	}

	return
}

func init() {
	UsefulFunc["HasKeyValueTypes"] = func(globals *model.Globals) bool {
		return len(KeyValueTypeNames(globals)) > 0
	}

	UsefulFunc["GetKeyValueTypeNames"] = KeyValueTypeNames

	UsefulFunc["GetIndicesByTable"] = GetIndicesByTable

	UsefulFunc["GetIndices"] = GetIndices
}
