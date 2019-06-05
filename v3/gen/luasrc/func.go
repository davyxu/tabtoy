package luasrc

import (
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func WrapValue(globals *model.Globals, value string, valueType *model.TypeDefine) string {
	if valueType.IsArray() {

		var sb strings.Builder
		sb.WriteString("{")

		// 空的单元格，导出空数组，除非强制指定填充默认值
		if value != "" {
			for index, elementValue := range strings.Split(value, valueType.ArraySplitter) {
				if index > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(gen.WrapSingleValue(globals, valueType, elementValue))
			}
		}

		sb.WriteString("}")

		return sb.String()

	} else {
		return gen.WrapSingleValue(globals, valueType, value)
	}

	return value
}

func init() {
	UsefulFunc["WrapTabValue"] = func(globals *model.Globals, dataTable *model.DataTable, allHeaders []*model.TypeDefine, row, col int) (ret string) {

		// 找到完整的表头（按完整表头遍历）
		header := allHeaders[col]

		if header == nil {
			return ""
		}

		// 在单元格找到值
		valueCell := dataTable.GetCell(row, col)

		if valueCell != nil {

			return WrapValue(globals, valueCell.Value, header)
		} else {
			// 这个表中没有这列数据
			return WrapValue(globals, "", header)
		}
	}

}
