package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
)

func matchField(objType reflect.Type, header string) int {

	for i := 0; i < objType.NumField(); i++ {
		fd := objType.Field(i)

		if fd.Tag.Get("tab_name") == header {
			return i
		}
	}

	return -1

}

// 将一行数据解析为具体的类型
func resolveRowType(ret *table.TypeField, tab *model.DataTable, row model.DataRow) {

	vobjtype := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for col, header := range tab.RawHeader() {

		index := matchField(tobj, header)

		if index == -1 {
			panic("类型表头找不到")
		}

		fieldType := vobjtype.Field(index)

		StringToValue(row[col], fieldType.Addr().Interface())
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

		resolveRowType(&objtype, tab, oneRow)

		globals.Symbols.AddField(&objtype)
	}

	return nil
}
