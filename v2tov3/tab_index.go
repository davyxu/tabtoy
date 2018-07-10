package v2tov3

import (
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
)

func ExportIndexTable(globals *model.Globals) error {

	var err error
	globals.TargetIndexSheet, err = globals.AddTable("Index.xlsx", "").AddSheet("Default")
	if err != nil {
		return err
	}
	helper.WriteIndexTableHeader(globals.TargetIndexSheet)

	globals.TargetTables.VisitAllTable(func(data *helper.MemFileData) bool {

		if data.FileName == "Index.xlsx" {
			return true
		}

		var mode string
		if data.FileName == "Type.xlsx" {
			mode = "类型表"
		} else {
			mode = "数据表"
		}

		helper.WriteRowValues(globals.TargetIndexSheet, mode, data.TableName+"Define", markFileNameUpgrade(data.FileName)+data.FileName)

		return true
	})

	return nil
}
