package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
)

// 将一行数据解析为具体的类型
func ParseRow(ret interface{}, tab *model.DataTable, row int) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for _, header := range tab.RawHeader() {

		value, tf := tab.GetValueByName(row, header)

		index := matchField(tobj, header)

		if index == -1 {
			panic("表头数据不匹配" + header)
		}

		fieldValue := vobj.Field(index)

		StringToValue(value, fieldValue.Addr().Interface(), tf)
	}
}

func loadPragma(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	tab, err := LoadTableData(fileName, nil)

	if err != nil {
		return err
	}

	ResolveHeaderFields(tab, "TablePragma", &globals.Symbols)

	var pragmaList []*table.TablePragma
	for row := 0; row < tab.RowCount(); row++ {

		var pragma table.TablePragma
		ParseRow(&pragma, tab, row)

		pragmaList = append(pragmaList, &pragma)
	}

	for _, pragma := range pragmaList {

		var tab *model.DataTable

		for _, tabName := range pragma.TableName {

			tab, err = LoadTableData(tabName+".xlsx", tab)

			if err != nil {
				return err
			}
		}

		ResolveHeaderFields(tab, tab.Name(), &globals.Symbols)

		globals.AddData(tab)

	}

	return nil
}
