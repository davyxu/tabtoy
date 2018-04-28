package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"path/filepath"
	"reflect"
	"strings"
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

func loadTocData(tab *model.DataTable, symbols *model.SymbolTable) (pragmaList []*table.TablePragma) {

	for row := 0; row < tab.RowCount(); row++ {

		var pragma table.TablePragma
		ParseRow(&pragma, tab, row, symbols)

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func loadToc(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	var tocTable = model.NewDataTable()
	err := LoadTableData(fileName, tocTable)

	if err != nil {
		return err
	}

	ResolveHeaderFields(tocTable, "TablePragma", globals.Symbols)

	pragmaList := loadTocData(tocTable, globals.Symbols)

	for _, pragma := range pragmaList {

		switch pragma.TableType {
		case table.TableType_DataTable:

			var tabName string

			// 表名空时，从文件名推断
			if pragma.TableName == "" {
				tabName = getFileName(pragma.TableFileName)
			} else {
				tabName = pragma.TableName
			}

			dataTable := globals.GetData(tabName)

			if dataTable == nil {
				dataTable = model.NewDataTable()
				dataTable.Name = tabName
				globals.AddData(dataTable)
			}

			err = LoadTableData(pragma.TableFileName, dataTable)

			if err != nil {
				return err
			}

		case table.TableType_SymbolTable:

			err = loadSymbols(globals, pragma.TableFileName)

			if err != nil {
				return err
			}

		}

	}

	for _, tab := range globals.Datas {
		ResolveHeaderFields(tab, tab.Name, globals.Symbols)
	}

	return nil
}

func getFileName(filename string) string {
	_, name := filepath.Split(filename)

	return strings.TrimSuffix(name, filepath.Ext(filename))
}
