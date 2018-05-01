package v2

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2/model"
)

func structFieldHasDefaultValue(structFD *model.FieldDescriptor) bool {

	d := structFD.Complex

	if d == nil {
		return false
	}

	for _, childFD := range d.Fields {

		if childFD.Meta.GetString("Default") != "" {
			return true
		}

	}

	return false
}

func mergeValues(modelData *model.DataModel, tab *model.Table, checker model.GlobalChecker) bool {

	var currFV *model.FieldValue

	for _, line := range modelData.Lines {

		record := model.NewRecord()

		for _, fv := range line.Values {

			currFV = fv

			var sugguestIgnore bool
			// repeated的, 没有填充的, 直接跳过, 不生成数据
			if fv.RawValue == "" && fv.FieldDef.Meta.GetString("Default") == "" {

				if !mustFillCheck(fv.FieldDef, fv.RawValue) {
					goto ErrorStop
				}

				if fv.FieldDef.IsRepeated {

					if fv.FieldDef.Type == model.FieldType_Struct {

						// 重复的 结构体字段, 且结构体字段没有默认值, 整个不导出
						if !structFieldHasDefaultValue(fv.FieldDef) {
							continue
						}

					} else {
						// 重复的普通字段导出, 做占位

					}

				} else {

					if fv.FieldDef.Type == model.FieldType_Struct {

						// 不重复的 结构体字段, 且结构体字段没有默认值, 整个不导出
						if !structFieldHasDefaultValue(fv.FieldDef) {

							sugguestIgnore = true
						}

					} else {

						// 非重复的普通字段不导出
						sugguestIgnore = true
					}

				}

			}

			if !coloumnProcessor(checker, record, fv.FieldDef, fv.RawValue, sugguestIgnore) {
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

	log.Errorf("%s|%s(%s)", currFV.FileName, currFV.SheetName, util.R1C1ToA1(currFV.R, currFV.C))
	return false

}
