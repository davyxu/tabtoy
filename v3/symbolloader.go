package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
)

// 将一行数据解析为具体的类型
func resolveRowType(rowValue interface{}, row model.DataRow) {

	vobjtype := reflect.ValueOf(rowValue).Elem()

	for col, value := range row {

		fieldType := vobjtype.Field(col)

		StringToValue(value, fieldType.Addr().Interface())
	}
}

func loadSymbols(globals *model.Globals, fileName string) error {

	tab, err := loadTable(fileName)

	if err != nil {
		return err
	}

	for row := 0; row < tab.RowCount(); row++ {

		var objtype table.TypeField

		oneRow := tab.GetDataRow(row)

		resolveRowType(&objtype, oneRow)

		globals.Symbols.AddField(&objtype)
	}

	return nil
}
