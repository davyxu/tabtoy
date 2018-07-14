package compiler

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
)

func LoadTypeTable(typeTab *model.TypeTable, indexGetter helper.FileGetter, fileName string, builtin bool) error {

	tabs, err := LoadDataTable(indexGetter, fileName, "TypeDefine")

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		ResolveHeaderFields(tab, "TypeDefine", typeTab)

		for row := 1; row < len(tab.Rows); row++ {

			var objtype model.TypeDefine

			if !ParseRow(&objtype, tab, row, typeTab) {
				continue
			}

			objtype.IsBuiltin = builtin

			if typeTab.FieldByName(objtype.ObjectType, objtype.FieldName) != nil {

				cell := tab.GetValueByName(row, "字段名")

				if cell != nil {
					report.ReportError("DuplicateTypeFieldName", cell.String(), objtype.ObjectType, objtype.FieldName)
				} else {
					report.ReportError("InvalidTypeTable", objtype.ObjectType, objtype.FieldName, tab.FileName)
				}

			}

			typeTab.AddField(&objtype, tab, row)
		}

	}

	return nil
}

func typeTable_CheckEnumValueEmpty(typeTab *model.TypeTable) {
	linq.From(typeTab.Raw()).WhereT(func(td *model.TypeData) bool {

		return td.Define.Kind == model.TypeUsage_Enum && td.Define.Value == ""
	}).ForEachT(func(td *model.TypeData) {

		cell := td.Tab.GetValueByName(td.Row, "值")

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

		if td.Define.IsBuiltin || td.Define.Kind != model.TypeUsage_Enum {
			continue
		}

		key := NameValuePair{td.Define.ObjectType, td.Define.Value}

		if _, ok := checker[key]; ok {

			cell := td.Tab.GetValueByName(td.Row, "值")

			report.ReportError("DuplicateEnumValue", cell.String())
		}

		checker[key] = td
	}
}

func CheckTypeTable(typeTab *model.TypeTable) {

	typeTable_CheckEnumValueEmpty(typeTab)

	typeTable_CheckDuplicateEnumValue(typeTab)
}
