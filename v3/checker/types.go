package checker

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func CheckTypes(tab *model.DataTable, types *model.TypeTable) {

	for col, headerType := range tab.HeaderFields {

		// 原始类型检查
		if !table.PrimitiveExists(headerType.FieldType) &&
			!types.ObjectExists(headerType.FieldType) { // 对象检查

			raw := tab.RawHeader[col]

			helper.ReportError("UnknownFieldType", raw.String())
		}
	}

}
