package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
)

func LoadTypeTable(typeTab *model.TypeTable, indexGetter util.FileGetter, fileName string) error {

	tabs, err := LoadDataTable(indexGetter, fileName, "TypeDefine", "TypeDefine", typeTab)

	if err != nil {
		return err
	}

	for _, tab := range tabs {

		for row := 1; row < len(tab.Rows); row++ {

			var objtype model.TypeDefine

			if !ParseRow(&objtype, tab, row, typeTab) {
				continue
			}

			if objtype.Kind == model.TypeUsage_None {
				util.ReportError("UnknownTypeKind", objtype.ObjectType, objtype.FieldName)
			}

			if typeTab.FieldByName(objtype.ObjectType, objtype.FieldName) != nil {
				cell := tab.GetValueByName(row, "字段名")
				if cell != nil {
					util.ReportError("DuplicateTypeFieldName", cell.String(), objtype.ObjectType, objtype.FieldName)
				} else {
					util.ReportError("InvalidTypeTable", objtype.ObjectType, objtype.FieldName, tab.FileName)
				}

			}

			typeTab.AddField(&objtype, tab, row)
		}

	}

	return nil
}
