package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/tealeg/xlsx"
)

type Sheet struct {
	*xlsx.Sheet

	Row int // 当前行

	Column int // 当前列

	file *File // 指向父级

	FieldHeader []*model.FieldDefine // 有效的字段行，可以做多sheet对比
}

// 取行列信息
func (self *Sheet) GetRC() (int, int) {

	return self.Row + 1, self.Column + 1

}

// 获取单元格
func (self *Sheet) GetCellData(cursor, index int) string {

	return strings.TrimSpace(self.Cell(cursor, index).Value)
}

// 设置单元格
func (self *Sheet) SetCellData(cursor, index int, data string) {

	self.Cell(cursor, index).Value = data
}

func NewSheet(file *File, sheet *xlsx.Sheet) *Sheet {
	self := &Sheet{
		file:  file,
		Sheet: sheet,
	}

	return self
}
