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

		index := matchField(tobj, header.Value)

		if index == -1 {
			ReportError("HeaderNotMatchFieldName", header, header.String())
		}

		fieldValue := vobj.Field(index)

		RawStringToValue(oneRow[col].Value, fieldValue.Addr().Interface())
	}

}

// 将一行数据解析为具体的类型
func ParseRow(ret interface{}, tab *model.DataTable, row int, symbols *model.SymbolTable) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for _, header := range tab.RawHeader {

		cell, tf := tab.GetValueByName(row, header.Value)

		index := matchField(tobj, header.Value)

		if index == -1 {
			ReportError("HeaderNotMatchFieldName", header.String())
		}

		fieldValue := vobj.Field(index)

		StringToValue(cell.Value, fieldValue.Addr().Interface(), tf, symbols)
	}
}
