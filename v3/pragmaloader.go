package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
)

// 将一行数据解析为具体的类型
func ParseRow(ret interface{}, tab *model.DataTable, row int, symbols *model.SymbolTable) {

	vobj := reflect.ValueOf(ret).Elem()

	tobj := reflect.TypeOf(ret).Elem()

	for _, header := range tab.RawHeader {

		value, tf := tab.GetValueByName(row, header)

		index := matchField(tobj, header)

		if index == -1 {
			panic("表头数据不匹配" + header)
		}

		fieldValue := vobj.Field(index)

		StringToValue(value, fieldValue.Addr().Interface(), tf, symbols)
	}
}

func loadPragmaData(tab *model.DataTable, symbols *model.SymbolTable) (pragmaList []*table.TablePragma) {

	for row := 0; row < tab.RowCount(); row++ {

		var pragma table.TablePragma
		ParseRow(&pragma, tab, row, symbols)

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func loadPragma(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	var pragmaTable = model.NewDataTable()
	err := LoadTableData(fileName, pragmaTable)

	if err != nil {
		return err
	}

	ResolveHeaderFields(pragmaTable, "TablePragma", globals.Symbols)

	pragmaList := loadPragmaData(pragmaTable, globals.Symbols)

	for _, pragma := range pragmaList {

		switch pragma.TableType {
		case table.TableType_DataTable:
			var dataTable = model.NewDataTable()
			dataTable.Name = pragma.TableName

			for _, fileName := range pragma.TableFileName {

				err = LoadTableData(fileName, dataTable)

				if err != nil {
					return err
				}
			}

			ResolveHeaderFields(dataTable, dataTable.Name, globals.Symbols)

			globals.AddData(dataTable)

		case table.TableType_SymbolTable:

			for _, fileName := range pragma.TableFileName {

				err = loadSymbols(globals, fileName)

				if err != nil {
					return err
				}
			}

		}

	}

	return nil
}
