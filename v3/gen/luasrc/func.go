package luasrc

import (
	"log"
	"strings"
	"text/template"

	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
)

var UsefulFunc = template.FuncMap{}

func WrapValue(globals *model.Globals, cell *model.Cell, valueType *model.TypeDefine) string {
	var fields = model.HadderStructCache[valueType.FieldName]
	if valueType.IsArray() {
		var sb strings.Builder
		sb.WriteString("{")

		if cell != nil {
			for index, elementValue := range cell.ValueList {
				if index > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(gen.WrapSingleValue(globals, valueType, elementValue))
			}
		}

		sb.WriteString("}")

		return sb.String()

	} else if fields != nil && len(fields.TypeInfo) > 0 {
		if len(cell.ValueList) < 1 {
			cell.ValueList = strings.Split(cell.Value, " ")
		}

		var sb strings.Builder
		sb.WriteString("{")

		if cell != nil {
			for index, elementValue := range cell.ValueList {
				if index > 0 {
					sb.WriteString(",")
				}
				data := strings.Split(elementValue, ":")
				if len(data) != 2 {
					log.Println(data, cell)
					report.ReportError("UnknownTypeKind", valueType.ObjectType, valueType.FieldName)
				}
				if _, ok := fields.TypeInfo[data[0]]; ok {
					sb.WriteString(fields.TypeInfo[data[0]].FieldName + " = ")
				} else {
					sb.WriteString(data[0] + " = ")
				}

				sb.WriteString(gen.WrapSingleValue(globals, valueType, data[1]))
			}
		}

		sb.WriteString("}")

		return sb.String()
	} else {

		var value string
		if cell != nil {
			value = cell.Value
		}

		return gen.WrapSingleValue(globals, valueType, value)
	}
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

			return WrapValue(globals, valueCell, header)
		} else {
			// 这个表中没有这列数据
			return WrapValue(globals, nil, header)
		}
	}

	UsefulFunc["IsWrapFieldName"] = func(globals *model.Globals, dataTable *model.DataTable, allHeaders []*model.TypeDefine, row, col int) (ret bool) {
		// 找到完整的表头（按完整表头遍历）
		header := allHeaders[col]

		if header == nil {
			return false
		}

		if globals.CanDoAction(model.ActionNoGennFieldLua, header) {
			return false
		}

		return true
	}

}
