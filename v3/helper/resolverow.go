package helper

import (
	"github.com/davyxu/tabtoy/v3/model"
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
func ResolveRowByReflect(ret interface{}, tab *model.DataTable, row int) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	oneRow := tab.GetDataRow(row)

	for col, header := range tab.RawHeader {

		index := matchField(tobj, header)

		if index == -1 {
			ReportError("HeaderNotMatchFieldName", header, Location(tab.FileName, 0, col))
		}

		fieldValue := vobj.Field(index)

		RawStringToValue(oneRow[col], fieldValue.Addr().Interface())
	}

}

// 将一行数据解析为具体的类型
func ParseRow(ret interface{}, tab *model.DataTable, row int, symbols *model.SymbolTable) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for col, header := range tab.RawHeader {

		value, tf := tab.GetValueByName(row, header)

		index := matchField(tobj, header)

		if index == -1 {
			ReportError("HeaderNotMatchFieldName", header, Location(tab.FileName, 0, col))
		}

		fieldValue := vobj.Field(index)

		StringToValue(value, fieldValue.Addr().Interface(), tf, symbols)
	}
}
