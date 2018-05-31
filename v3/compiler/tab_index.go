package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"path/filepath"
	"strings"
)

func loadIndexData(tab *model.DataTable, symbols *model.TypeTable) (pragmaList []*table.TablePragma) {

	for row := 0; row < len(tab.Rows); row++ {

		var pragma table.TablePragma
		model.ParseRow(&pragma, tab, row, symbols)

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func LoadIndexTable(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	tabs, err := LoadDataTable(globals.IndexGetter, fileName, "TablePragma")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TablePragma", globals.Types)

		globals.IndexList = loadIndexData(tab, globals.Types)
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
