package v3

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"path/filepath"
	"strings"
)

func loadIndexData(tab *model.DataTable, symbols *model.SymbolTable) (pragmaList []*table.TablePragma) {

	for row := 0; row < tab.RowCount(); row++ {

		var pragma table.TablePragma
		helper.ParseRow(&pragma, tab, row, symbols)

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func LoadIndex(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	var indexTab = model.NewDataTable()
	indexTab.FileName = fileName
	indexTab.Name = "TablePragma"

	err := LoadTableData(fileName, indexTab)

	if err != nil {
		return err
	}

	ResolveHeaderFields(indexTab, "TablePragma", globals.Symbols)

	globals.IndexList = loadIndexData(indexTab, globals.Symbols)

	return nil
}

// 表名空时，从文件名推断
func getTableName(pragma *table.TablePragma) string {

	if pragma.TableType == "" {

		_, name := filepath.Split(pragma.TableFileName)

		return strings.TrimSuffix(name, filepath.Ext(pragma.TableFileName))
	} else {
		return pragma.TableType
	}
}
