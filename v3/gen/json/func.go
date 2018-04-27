package json

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func init() {
	UsefulFunc["WrapJsonValue"] = func(globals *model.Globals, dataTable *model.DataTable, row, col int) (ret string) {

		// 单元格的值
		value := dataTable.GetValue(row, col)

		// 表头的类型
		valueType := dataTable.GetType(col)

		switch {
		case valueType.FieldType == "string": // 字符串
			return util.StringEscape(value)
		case globals.Symbols.IsEnumKind(dataTable.Name(), valueType.FieldType): // 枚举
			return globals.Symbols.ResolveEnumValue(valueType.FieldType, value)
		}

		return value
	}
}
