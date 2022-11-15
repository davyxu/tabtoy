package util

import (
	"github.com/tealeg/xlsx"
	"strings"
)

type XlsxFile struct {
	file *xlsx.File

	sheets   []TableSheet
	cacheDir string
}

func (self *XlsxFile) Sheets() (ret []TableSheet) {

	return self.sheets
}

func (self *XlsxFile) Save(filename string) error {
	return self.file.Save(filename)
}

func (self *XlsxFile) Load(filename string) (err error) {

	var file *xlsx.File

	if self.cacheDir == "" {
		file, err = xlsx.OpenFile(filename)
		if err != nil {
			return err
		}
	} else {
		cache := NewTableCache(filename, self.cacheDir)

		if err = cache.Open(); err != nil {
			return err
		}

		if file, err = cache.Load(); err != nil {
			return err
		} else {

			if !cache.UseCache() {
				cache.Save()
			}
		}
	}

	self.FromXFile(file)

	return nil
}

func (self *XlsxFile) FromXFile(file *xlsx.File) {
	self.file = file

	for _, sheet := range file.Sheets {
		self.sheets = append(self.sheets, newXlsxSheet(sheet))
	}
}

func NewXlsxFile(cacheDir string) TableFile {

	self := &XlsxFile{
		cacheDir: cacheDir,
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

func (self *XlsxSheet) IsRowEmpty(row, maxCol int) bool {

	if maxCol == -1 {
		maxCol = self.Sheet.MaxCol
	}

	for col := 0; col < maxCol; col++ {

		data := self.GetValue(row, col, nil)

		if data != "" {
			return false
		}
	}

	return true
}

func (self *XlsxSheet) GetValue(row, col int, opt *ValueOption) (ret string) {
	c := self.Sheet.Cell(row, col)

	// 浮点数单元格按原样输出
	if opt != nil && opt.ValueAsFloat {
		ret, _ = c.GeneralNumeric()
		ret = strings.TrimSpace(ret)
	} else {
		// 取列头所在列和当前行交叉的单元格
		ret = strings.TrimSpace(c.Value)
	}

	return
}

func (self *XlsxSheet) WriteRow(valueList ...string) {
	row := self.Sheet.AddRow()

	for _, value := range valueList {

		cell := row.AddCell()
		cell.SetValue(value)
	}
}

func newXlsxSheet(sheet *xlsx.Sheet) TableSheet {
	return &XlsxSheet{
		Sheet: sheet,
	}
}
