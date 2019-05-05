package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"path/filepath"
	"sort"
	"strings"
)

func parseIndexRow(tab *model.DataTable, symbols *model.TypeTable) (pragmaList []*model.IndexDefine) {

	for row := 1; row < len(tab.Rows); row++ {

		var pragma model.IndexDefine
		if !ParseRow(&pragma, tab, row, symbols) {
			continue
		}

		if pragma.Kind == model.TableKind_Type {
			pragma.TableType = "TypeDefine"
		}

		if pragma.TableType == "" {

			_, name := filepath.Split(pragma.TableFileName)

			pragma.TableType = strings.TrimSuffix(name, filepath.Ext(pragma.TableFileName))
		}

		pragmaList = append(pragmaList, &pragma)
	}

	return
}

func LoadIndexTable(globals *model.Globals, fileName string) error {

	if fileName == "" {
		return nil
	}

	// 加载原始数据
	tabs, err := LoadDataTable(globals.IndexGetter, fileName, "IndexDefine", "IndexDefine", globals.Types)

	if err != nil {
		return err
	}

	var pragmaList []*model.IndexDefine

	for _, tab := range tabs {
		pragmaList = append(pragmaList, parseIndexRow(tab, globals.Types)...)
	}

	// 按表类型排序，保证类型表先读取
	sort.Slice(pragmaList, func(i, j int) bool {
		a := pragmaList[i]
		b := pragmaList[j]

		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}

		if a.TableType != b.TableType {
			return a.TableType < b.TableType
		}

		return a.TableFileName < b.TableFileName

	})

	globals.IndexList = pragmaList

	return nil
}
