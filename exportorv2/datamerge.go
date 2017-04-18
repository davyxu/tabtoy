package exportorv2

import (
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

func mergeValues(modelData *model.DataModel, tab *model.Table, checker model.GlobalChecker, vertical bool) bool {

	var currFV *model.FieldValue

	for _, line := range modelData.Lines {

		record := model.NewRecord()

		for _, fv := range line.Values {

			// repeated的, 没有填充的, 直接跳过, 不生成数据
			if fv.RawValue == "" && fv.FieldDef.Meta.GetString("Default") == "" {

				if !mustFillCheck(fv.FieldDef, fv.RawValue) {
					goto ErrorStop
				}

				// 空的, 不重复字段忽略
				if !fv.FieldDef.IsRepeated {
					continue
				}

				// 空的, 重复的, 结构体字段忽略
				if fv.FieldDef.Type == model.FieldType_Struct {
					continue
				}

				// 只保留重复的普通字段

			}

			if !coloumnProcessor(checker, record, fv.FieldDef, fv.RawValue) {
				goto ErrorStop
			}

		}

		tab.Add(record)
	}

	return true

ErrorStop:

	if currFV == nil {
		return false
	}

	log.Errorf("%s|%s(%s)", currFV.FileName, currFV.SheetName, util.ConvR1C1toA1(currFV.R, currFV.C))
	return false

}
