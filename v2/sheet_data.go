package v2

import (
	"strings"

	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

/*
	Sheet数据表单的处理

*/

type DataSheet struct {
	*Sheet
}

func (self *DataSheet) Valid() bool {

	name := strings.TrimSpace(self.Sheet.Name)
	if name != "" && name[0] == '#' {
		return false
	}

	return self.GetCellData(0, 0) != ""
}

func (self *DataSheet) Export(file *File, dataModel *model.DataModel, dataHeader, parentHeader *DataHeader) bool {

	verticalHeader := file.LocalFD.Pragma.GetBool("Vertical")

	if verticalHeader {
		return self.exportColumnMajor(file, dataModel, dataHeader, parentHeader)
	} else {
		return self.exportRowMajor(file, dataModel, dataHeader, parentHeader)
	}

}

// 导出以行数据延展的表格(普通表格)
func (self *DataSheet) exportRowMajor(file *File, dataModel *model.DataModel, dataHeader, parentHeader *DataHeader) bool {

	// 是否继续读行
	var readingLine bool = true

	var meetEmptyLine bool

	var warningAfterEmptyLineDataOnce bool

	// 遍历每一行
	for self.Row = DataSheetHeader_DataBegin; readingLine; self.Row++ {

		// 整行都是空的
		if self.IsFullRowEmpty(self.Row, dataHeader.RawFieldCount()) {

			// 再次碰空行, 表示确实是空的
			if meetEmptyLine {
				break

			} else {
				meetEmptyLine = true
			}

			continue

		} else {

			//已经碰过空行, 这里又碰到数据, 说明有人为隔出的空行, 做warning提醒, 防止数据没导出
			if meetEmptyLine && !warningAfterEmptyLineDataOnce {
				r, _ := self.GetRC()

				log.Warnf("%s %s|%s(%s)", i18n.String(i18n.DataSheet_RowDataSplitedByEmptyLine), self.file.FileName, self.Name, util.ConvR1C1toA1(r, 1))

				warningAfterEmptyLineDataOnce = true
			}

			// 曾经有过空行, 即便现在不是空行也没用, 结束
			if meetEmptyLine {
				break
			}

		}

		line := model.NewLineData()

		// 遍历每一列
		for self.Column = 0; self.Column < dataHeader.RawFieldCount(); self.Column++ {

			fieldDef, ok := fieldDefGetter(self.Column, dataHeader, parentHeader)

			if !ok {
				log.Errorf("%s %s|%s(%s)", i18n.String(i18n.DataHeader_FieldNotDefinedInMainTableInMultiTableMode), self.file.FileName, self.Name, util.ConvR1C1toA1(self.Row+1, self.Column+1))
				return false
			}

			op := self.processLine(fieldDef, line, dataHeader)

			if op == lineOp_Continue {
				continue
			} else if op == lineOp_Break {
				break
			}

		}

		// 是子表
		if parentHeader != nil {

			// 遍历母表所有的列头字段
			for c := 0; c < parentHeader.RawFieldCount(); c++ {
				fieldDef := parentHeader.RawField(c)

				// 在子表中有对应字段的, 忽略, 只要没有的字段
				if _, ok := dataHeader.HeaderByName[fieldDef.Name]; ok {
					continue
				}

				op := self.processLine(fieldDef, line, dataHeader)

				if op == lineOp_Continue {
					continue
				} else if op == lineOp_Break {
					break
				}

			}
		}

		dataModel.Add(line)

	}

	return true
}

const (
	lineOp_none = iota
	lineOp_Break
	lineOp_Continue
)

func (self *DataSheet) processLine(fieldDef *model.FieldDescriptor, line *model.LineData, dataHeader *DataHeader) int {
	// 数据大于列头时, 结束这个列
	if fieldDef == nil {
		return lineOp_Break
	}

	// #开头表示注释, 跳过
	if strings.Index(fieldDef.Name, "#") == 0 {
		return lineOp_Continue
	}

	var rawValue string
	// 浮点数按本来的格式输出
	if fieldDef.Type == model.FieldType_Float && !fieldDef.IsRepeated {
		rawValue = self.GetCellDataAsNumeric(self.Row, self.Column)
	} else {
		rawValue = self.GetCellData(self.Row, self.Column)
	}

	r, c := self.GetRC()

	line.Add(&model.FieldValue{
		FieldDef:           fieldDef,
		RawValue:           rawValue,
		SheetName:          self.Name,
		FileName:           self.file.FileName,
		R:                  r,
		C:                  c,
		FieldRepeatedCount: dataHeader.FieldRepeatedCount(fieldDef),
	})

	return lineOp_none
}

// 多表合并时, 要从从表的字段名在主表的表头里做索引
func fieldDefGetter(index int, dataHeader, parentHeader *DataHeader) (*model.FieldDescriptor, bool) {

	fieldDef := dataHeader.RawField(index)
	if fieldDef == nil {
		return nil, true
	}

	if parentHeader != nil {

		if strings.Index(fieldDef.Name, "#") == 0 {
			return fieldDef, true
		}

		ret, ok := parentHeader.HeaderByName[fieldDef.Name]
		if !ok {
			return nil, false
		}
		return ret, true
	}

	return fieldDef, true

}

func mustFillCheck(fd *model.FieldDescriptor, raw string) bool {
	// 值重复检查
	if fd.Meta.GetBool("MustFill") {

		if raw == "" {
			log.Errorf("%s, %s", i18n.String(i18n.DataSheet_MustFill), fd.String())
			return false
		}
	}

	return true
}

func newDataSheet(sheet *Sheet) *DataSheet {

	return &DataSheet{
		Sheet: sheet,
	}
}
