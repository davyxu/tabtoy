package v3

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/davyxu/tabtoy/v3/table"
)

var coreSymbols model.TypeTable

func init() {

	for _, symbol := range table.CoreSymbols {
		coreSymbols.AddField(symbol, nil, 0)
	}
}

func LoadTypeTable(typeTab *model.TypeTable, indexGetter helper.FileGetter, fileName string, builtin bool) error {

	tabs, err := LoadDataTable(indexGetter, fileName, "TableField")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TableField", &coreSymbols)

		for row := range tab.Rows {

			var objtype table.TableField

			model.ParseRow(&objtype, tab, row, &coreSymbols)

			objtype.IsBuiltin = builtin

			if typeTab.FieldByName(objtype.ObjectType, objtype.FieldName) != nil {

				cell, _ := tab.GetValueByName(row, "字段名")

				report.ReportError("DuplicateTypeFieldName", cell.String())
			}

			typeTab.AddField(&objtype, tab, row)
		}

	}

	return nil
}

func typeTable_CheckEnumValueEmpty(typeTab *model.TypeTable) {
	linq.From(typeTab.Raw()).WhereT(func(td *model.TypeData) bool {

		return td.Type.Kind == table.TableKind_Enum && td.Type.Value == ""
	}).ForEachT(func(td *model.TypeData) {

		cell, _ := td.Tab.GetValueByName(td.Row, "值")

		report.ReportError("EnumValueEmpty", cell.String())
	})
}

func typeTable_CheckDuplicateEnumValue(typeTab *model.TypeTable) {

	type NameValuePair struct {
		Name  string
		Value string
	}

	checker := map[NameValuePair]*model.TypeData{}

	for _, td := range typeTab.Raw() {

		if td.Type.IsBuiltin || td.Type.Kind != table.TableKind_Enum {
			continue
		}

		key := NameValuePair{td.Type.ObjectType, td.Type.Value}

		if _, ok := checker[key]; ok {

			cell, _ := td.Tab.GetValueByName(td.Row, "值")

			report.ReportError("DuplicateEnumValue", cell.String())
		}

		checker[key] = td
	}
}

func checkValueType(typeDataTab *model.DataTable) {

	//for row := range typeDataTab.Rows {
	//
	//	for col := range typeDataTab.RawHeader {
	//
	//		value := typeDataTab.GetValue(row, col)
	//
	//		valueType := typeDataTab.GetType(col)
	//
	//	}
	//
	//}

}

func CheckTypeTable(typeTab *model.TypeTable) {

	typeTable_CheckEnumValueEmpty(typeTab)

	typeTable_CheckDuplicateEnumValue(typeTab)
}
