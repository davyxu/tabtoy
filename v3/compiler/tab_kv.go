package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
)

func transposeKVtoData(symbols *model.TypeTable, kvtab *model.DataTable) (ret *model.DataTable) {

	ret = model.NewDataTable()
	ret.HeaderType = kvtab.HeaderType
	ret.OriginalHeaderType = kvtab.HeaderType
	ret.FileName = kvtab.FileName
	ret.SheetName = kvtab.SheetName

	// 添加表头
	ret.AddRow()

	// 添加数据行
	ret.AddRow()

	// 遍历KV表的每一行
	for row := 1; row < len(kvtab.Rows); row++ {

		fieldName := kvtab.GetValueByName(row, "字段名")
		fieldType := kvtab.GetValueByName(row, "字段类型")
		name := kvtab.GetValueByName(row, "标识名")

		arraySplitter := kvtab.GetValueByName(row, "数组切割")

		var tf table.TableField
		tf.Kind = table.TableKind_HeaderStruct
		tf.ObjectType = kvtab.HeaderType

		tf.Name = name.Value

		tf.FieldName = fieldName.Value
		tf.FieldType = fieldType.Value
		tf.ArraySplitter = arraySplitter.Value

		if symbols.FieldByName(tf.ObjectType, tf.FieldName) != nil {
			report.ReportError("DuplicateKVField", fieldName.String())
		}

		symbols.AddField(&tf, kvtab, row)

		// 输出表的表头原始数据
		headerCell := ret.AddCell(0)
		headerCell.Value = fieldName.Value

		header := ret.MustGetHeader(headerCell.Col)
		header.Cell.Value = fieldName.Value
		header.TypeInfo = &tf

		inputValueCell := kvtab.GetValueByName(row, "值")

		outputValueCell := ret.AddCell(1)
		outputValueCell.CopyFrom(inputValueCell)

	}

	return
}
