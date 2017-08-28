package v2

import (
	"strings"

	"github.com/tealeg/xlsx"
)

// 描述一个表单
type Sheet struct {
	*xlsx.Sheet

	Row int // 当前行

	Column int // 当前列

	file *File // 指向父级

}

// 取行列信息
func (self *Sheet) GetRC() (int, int) {

	return self.Row + 1, self.Column + 1

}

// 获取单元格 cursor=行,  index=列
func (self *Sheet) GetCellData(cursor, index int) string {

	if cursor >= len(self.Rows) {
		return ""
	}

	r := self.Rows[cursor]
	for len(r.Cells) <= index {
		r.AddCell()
	}

	return strings.TrimSpace(r.Cells[index].Value)
}

// 设置单元格
func (self *Sheet) SetCellData(cursor, index int, data string) {

	self.Cell(cursor, index).Value = data
}

// 整行都是空的
func (self *Sheet) IsFullRowEmpty(row, maxCol int) bool {

	for col := 0; col < maxCol; col++ {

		data := self.GetCellData(row, col)

		if data != "" {
			return false
		}
	}

	return true
}

func NewSheet(file *File, sheet *xlsx.Sheet) *Sheet {
	self := &Sheet{
		file:  file,
		Sheet: sheet,
	}

	return self
}
