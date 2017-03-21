package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

const (
	ColumnMajor_RowDataBegin = 1
	ColumnMajor_ColumnValue  = 4
)

// 导出适合配置格式的表格
func (self *DataSheet) exportColumnMajor(file *File, dataModel *model.DataModel, dataHeader, parentHeader *DataHeader) bool {

	// 是否继续读行
	var readingLine bool = true

	var meetEmptyLine bool

	var warningAfterEmptyLineDataOnce bool

	line := model.NewLineData()

	for self.Row = ColumnMajor_RowDataBegin; readingLine; self.Row++ {
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

		}

		fieldDef := dataHeader.RawField(self.Row - ColumnMajor_RowDataBegin)

		// 数据大于列头时, 结束这个列
		if fieldDef == nil {
			break
		}

		// #开头表示注释, 跳过
		if strings.Index(fieldDef.Name, "#") == 0 {
			continue
		}

		rawValue := self.GetCellData(self.Row, ColumnMajor_ColumnValue)

		r, c := self.GetRC()

		line.Add(&model.FieldValue{
			FieldDef:  fieldDef,
			RawValue:  rawValue,
			SheetName: self.Name,
			R:         r,
			C:         c,
			File:      file,
		})

	}

	dataModel.Add(line)

	return true

}
