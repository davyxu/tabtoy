package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func transposeKVtoData(symbols *model.SymbolTable, kvtab *model.DataTable) (ret *model.DataTable) {

	ret = model.NewDataTable()
	ret.HeaderType = kvtab.HeaderType
	ret.OriginalHeaderType = kvtab.HeaderType
	ret.FileName = kvtab.FileName

	var oneRow model.DataRow
	for row := range kvtab.Rows {

		fieldName, _ := kvtab.GetValueByName(row, "字段名")
		fieldType, _ := kvtab.GetValueByName(row, "字段类型")
		name, _ := kvtab.GetValueByName(row, "标识名")
		arraySplitter, _ := kvtab.GetValueByName(row, "数组切割")

		var tf table.TableField
		tf.Kind = "表头"
		tf.ObjectType = kvtab.HeaderType

		tf.Name = name

		tf.FieldName = fieldName
		tf.FieldType = fieldType
		tf.ArraySplitter = arraySplitter

		if symbols.FindField(tf.ObjectType, tf.FieldName) != nil {
			helper.ReportError("DuplicateKVField", helper.Location(kvtab.FileName, row, 0))
		}

		symbols.AddField(&tf)

		value, _ := kvtab.GetValueByName(row, "值")

		oneRow = append(oneRow, value)

		ret.AddHeaderField(&tf)
	}

	// KV只有一行，列是原表的行
	ret.AddRow(oneRow)

	return
}
