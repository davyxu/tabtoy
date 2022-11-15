package checker

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"go/token"
)

func CheckType(typeTab *model.TypeTable) {

	typeTable_CheckField(typeTab)

	typeTable_CheckEnumValueEmpty(typeTab)

	typeTable_CheckDuplicateEnumValue(typeTab)

}

func isValidFieldName(name string) bool {

	return token.IsIdentifier(name)
}

func typeTable_CheckField(typeTab *model.TypeTable) {
	for _, td := range typeTab.Raw() {

		if !isValidFieldName(td.Define.FieldName) {
			cell := td.Tab.GetValueByName(td.Row, "字段名")
			util.ReportError("InvalidFieldName", cell.String())
		}
	}
}

func typeTable_CheckEnumValueEmpty(typeTab *model.TypeTable) {
	linq.From(typeTab.Raw()).Where(func(raw interface{}) bool {
		td := raw.(*model.TypeData)

		return td.Define.Kind == model.TypeUsage_Enum && td.Define.Value == ""
	}).ForEach(func(raw interface{}) {
		td := raw.(*model.TypeData)

		cell := td.Tab.GetValueByName(td.Row, "值")

		util.ReportError("EnumValueEmpty", cell.String())
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

			util.ReportError("DuplicateEnumValue", cell.String())
		}

		checker[key] = td
	}
}
