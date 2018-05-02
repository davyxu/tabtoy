package json

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func wrapSingleValue(globals *model.Globals, valueType *table.TableField, value string) string {
	switch {
	case valueType.FieldType == "string": // 字符串
		return util.StringEscape(value)
	case globals.Symbols.IsEnumKind(valueType.FieldType): // 枚举
		return globals.Symbols.ResolveEnumValue(valueType.FieldType, value)
	case valueType.FieldType == "bool":

		switch value {
		case "是", "yes", "YES", "1", "true", "TRUE", "True":
			return "true"
		case "否", "no", "NO", "0", "false", "FALSE", "False":
			return "false"
		}

		return "false"
	}

	if value == "" {
		return table.FetchDefaultValue(valueType)
	}

	return value

}

func WrapValue(globals *model.Globals, value string, valueType *table.TableField) string {
	if valueType.IsArray() {

		var sb strings.Builder
		sb.WriteString("[")

		// 空的单元格，导出空数组，除非强制指定填充默认值
		if value != "" {
			for index, elementValue := range strings.Split(value, valueType.ArraySplitter) {
				if index > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(wrapSingleValue(globals, valueType, elementValue))
			}
		}

		sb.WriteString("]")

		return sb.String()

	} else {
		return wrapSingleValue(globals, valueType, value)
	}

	return value
}

func init() {
	UsefulFunc["WrapTabValue"] = func(globals *model.Globals, dataTable *model.DataTable, row, col int) (ret string) {

		// 单元格的值
		value := dataTable.GetValue(row, col)

		// 表头的类型
		valueType := dataTable.GetType(col)

		return WrapValue(globals, value, valueType)
	}

}
