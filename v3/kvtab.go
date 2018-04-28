package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func convertKVToData(symbols *model.SymbolTable, kvtab *model.DataTable) (ret *model.DataTable) {

	ret = model.NewDataTable()
	ret.Name = kvtab.Name

	var oneRow model.DataRow
	for row := range kvtab.Rows {

		fieldName, _ := kvtab.GetValueByName(row, "字段名")
		fieldType, _ := kvtab.GetValueByName(row, "字段类型")
		name, _ := kvtab.GetValueByName(row, "标识名")
		isArray, _ := kvtab.GetValueByName(row, "数组")
		splitter, _ := kvtab.GetValueByName(row, "切割符")

		var tf table.TableField
		tf.Kind = "表头"
		tf.ObjectType = kvtab.Name
		tf.Name = name
		tf.FieldName = fieldName
		tf.FieldType = fieldType
		RawStringToValue(isArray, &tf.IsArray)
		tf.Splitter = splitter
		symbols.AddField(&tf)

		value, _ := kvtab.GetValueByName(row, "值")

		oneRow = append(oneRow, value)

		ret.AddHeaderField(&tf)
	}

	// KV只有一行，列是原表的行
	ret.AddRow(oneRow)

	return
}
