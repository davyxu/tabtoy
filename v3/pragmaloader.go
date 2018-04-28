package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"path/filepath"
	"strings"
)

func loadIndexData(tab *model.DataTable, symbols *model.SymbolTable) (pragmaList []*table.TablePragma) {

	for row := 0; row < tab.RowCount(); row++ {

		var pragma table.TablePragma
		ParseRow(&pragma, tab, row, symbols)

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func LoadIndex(globals *model.Globals, fileName string, callback func(*table.TablePragma) error) error {

	if fileName == "" {
		return nil
	}

	var tocTable = model.NewDataTable()
	err := LoadTableData(fileName, tocTable)

	if err != nil {
		return err
	}

	ResolveHeaderFields(tocTable, "TablePragma", globals.Symbols)

	pragmaList := loadIndexData(tocTable, globals.Symbols)

	for _, pragma := range pragmaList {

		if err = callback(pragma); err != nil {
			return err
		}
	}

	return nil
}

// 表名空时，从文件名推断
func getTableName(pragma *table.TablePragma) string {

	if pragma.TableName == "" {

		_, name := filepath.Split(pragma.TableFileName)

		return strings.TrimSuffix(name, filepath.Ext(pragma.TableFileName))
	} else {
		return pragma.TableName
	}
}
