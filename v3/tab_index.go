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

func LoadIndexTable(globals *model.Globals, indexGetter FileGetter, fileName string) error {

	if fileName == "" {
		return nil
	}

	tabs, err := LoadDataTable(indexGetter, fileName, "TablePragma")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TablePragma", globals.Symbols)

		globals.IndexList = loadIndexData(tab, globals.Symbols)
	}

	return nil
}

// 表名空时，从文件名推断
func fillTableType(pragma *table.TablePragma) {

	if pragma.TableType == "" {

		_, name := filepath.Split(pragma.TableFileName)

		pragma.TableType = strings.TrimSuffix(name, filepath.Ext(pragma.TableFileName))
	}

}
