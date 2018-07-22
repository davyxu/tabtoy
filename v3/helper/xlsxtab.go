package helper

import "github.com/tealeg/xlsx"

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

func (self *XlsxSheet) IsFullRowEmpty(row int) bool {
	return IsFullRowEmpty(self.Sheet, row)
}

func (self *XlsxSheet) GetValue(row, col int, isFloat bool) string {

	// 浮点数单元格按原样输出
	if isFloat {
		return GetSheetValueAsNumericString(self.Sheet, row, col)
	} else {
		// 取列头所在列和当前行交叉的单元格
		return GetSheetValueString(self.Sheet, row, col)
	}
}

func NewXlsxSheet(sheet *xlsx.Sheet) TableSheet {
	return &XlsxSheet{
		Sheet: sheet,
	}
}
