package v2tov3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
)

func importDatas(globals *model.Globals, sourceSheet, targetSheet *xlsx.Sheet, headerList []model.ObjectFieldType) error {

	var row, col int

	for row = 4; ; row++ {

		if helper.IsFullRowEmpty(sourceSheet, row) {
			break
		}
		rowData := targetSheet.AddRow()

		var header model.ObjectFieldType

		for col, header = range headerList {

			sourceCell := sourceSheet.Cell(row, col)

			targetCell := rowData.AddCell()

			if header.IsArray() {
				targetCell.SetValue(sourceCell.Value)
				continue
			}

			setTargetCell(header.FieldType, sourceCell, targetCell, row, col)
		}

	}

	return nil

}

func setTargetCell(headerFieldType string, sourceCell, targetCell *xlsx.Cell, row, col int) {
	switch headerFieldType {
	case "int32", "uint32":

		if sourceCell.Value == "" {
			targetCell.SetInt(0)
			break
		}

		v, err := sourceCell.Int()
		if err != nil {
			log.Errorf("单元格转换错误 @%s, %s", util.R1C1ToA1(row+1, col+1), err.Error())
		} else {
			targetCell.SetInt(v)
		}

	case "int64", "uint64":
		v, err := sourceCell.Int64()
		if err != nil {
			log.Errorf("单元格转换错误 @%s, %s", util.R1C1ToA1(row+1, col+1), err.Error())
		} else {
			targetCell.SetInt64(v)
		}
	case "float":
		if sourceCell.Value == "" {
			targetCell.SetFloat(0)
			break
		}

		v, err := sourceCell.Float()
		if err != nil {
			log.Errorf("单元格转换错误 @%s, %s", util.R1C1ToA1(row+1, col+1), err.Error())
		} else {
			targetCell.SetFloat(v)
		}
	case "bool":
		var v bool
		if err, _ := util.StringToPrimitive(sourceCell.Value, &v); err != nil {
			log.Errorf("单元格转换错误 @%s, %s", util.R1C1ToA1(row+1, col+1), err.Error())
		} else {
			targetCell.SetBool(v)
		}
	default:
		targetCell.SetValue(sourceCell.Value)
	}
}
