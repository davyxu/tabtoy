package helper

import (
	"github.com/tealeg/xlsx"
	"strings"
)

type XlsxFile struct {
	*xlsx.File

	sheets []TableSheet
}

func (self *XlsxFile) Sheets() (ret []TableSheet) {

	return self.sheets
}

func NewXlsxFile(file *xlsx.File) TableFile {

	self := &XlsxFile{
		File: file,
	}

	for _, sheet := range file.Sheets {
		self.sheets = append(self.sheets, NewXlsxSheet(sheet))
	}

	return self
}

type XlsxSheet struct {
	*xlsx.Sheet
}

func (self *XlsxSheet) Name() string {
	return self.Sheet.Name
}

func (self *XlsxSheet) MaxColumn() int {
	return self.Sheet.MaxCol
}

func (self *XlsxSheet) IsFullRowEmpty(row int) bool {

	for col := 0; col < self.Sheet.MaxCol; col++ {

		data := self.GetValue(row, col, false)

		if data != "" {
			return false
		}
	}

	return true
}

func (self *XlsxSheet) GetValue(row, col int, isFloat bool) (ret string) {

	c := self.Sheet.Cell(row, col)

	// 浮点数单元格按原样输出
	if isFloat {
		ret, _ = c.GeneralNumeric()
		ret = strings.TrimSpace(ret)
	} else {
		// 取列头所在列和当前行交叉的单元格
		ret = strings.TrimSpace(c.Value)
	}

	return
}

func NewXlsxSheet(sheet *xlsx.Sheet) TableSheet {
	return &XlsxSheet{
		Sheet: sheet,
	}
}
