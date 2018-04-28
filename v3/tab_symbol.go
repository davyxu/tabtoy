package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
)

func matchField(objType reflect.Type, header string) int {

	for i := 0; i < objType.NumField(); i++ {
		fd := objType.Field(i)

		if fd.Tag.Get("tb_name") == header {
			return i
		}
	}

	return -1

}

// 将一行数据解析为具体的类型
func resolveRowTypeByReflect(ret interface{}, tab *model.DataTable, row int) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	oneRow := tab.GetDataRow(row)

	for col, header := range tab.RawHeader {

		index := matchField(tobj, header)

		if index == -1 {
			panic("表头数据不匹配:" + header)
		}

		fieldValue := vobj.Field(index)

		RawStringToValue(oneRow[col], fieldValue.Addr().Interface())
	}

}

func loadSymbols(globals *model.Globals, fileName string) error {

	var symbolTable = model.NewDataTable()
	err := LoadTableData(fileName, symbolTable)

	if err != nil {
		return err
	}

	for row := 0; row < symbolTable.RowCount(); row++ {

		var objtype table.TableField

		resolveRowTypeByReflect(&objtype, symbolTable, row)

		globals.Symbols.AddField(&objtype)
	}

	return nil
}
