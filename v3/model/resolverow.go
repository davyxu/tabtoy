package model

import (
	"github.com/davyxu/tabtoy/v3/report"
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
func ParseRow(ret interface{}, tab *DataTable, row int, symbols *TypeTable) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for _, header := range tab.Headers {

		valueCell := tab.GetValueByName(row, header.Cell.Value)

		if valueCell == nil {
			report.Log.Errorln("空")
			continue
		}

		index := matchField(tobj, header.Cell.Value)

		if index == -1 {
			report.ReportError("HeaderNotMatchFieldName", header.Cell.String())
		}

		fieldValue := vobj.Field(index)

		if err := StringToValue(valueCell.Value, fieldValue.Addr().Interface(), header.TypeInfo, symbols); err != nil {
			panic(err)
		}
	}
}
