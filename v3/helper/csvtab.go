package helper

import (
	"bytes"
	"encoding/csv"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type CSVFile struct {
	sheet *CSVSheet

	Name string

	records [][]string
}

func (self *CSVFile) Sheets() (ret []TableSheet) {

	return []TableSheet{self.sheet}
}

func (self *CSVFile) MaxCol() int {

	if len(self.records) > 0 {
		return len(self.records[0])
	}

	return 0
}

func ConvGBKToUTF8(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (self *CSVFile) ConvertToUTF8() {

	for row, rowData := range self.records {

		for col, cell := range rowData {

			if cell != "" {
				data, _ := ConvGBKToUTF8([]byte(cell))
				self.records[row][col] = string(data)
			}
		}

	}

}

func NewCSVFile(filename string) (*CSVFile, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	self := &CSVFile{}

	self.Name = strings.TrimSuffix(filename, filepath.Ext(filename))

	self.sheet = &CSVSheet{file: self}

	self.records, err = csv.NewReader(bytes.NewReader(data)).ReadAll()

	// 电子表格默认只能导入GBK的编码，因此认为输入的都是GBK的CSV
	self.ConvertToUTF8()

	return self, err
}

type CSVSheet struct {
	file *CSVFile
}

func (self *CSVSheet) Name() string {
	return self.file.Name
}

func (self *CSVSheet) IsFullRowEmpty(row int) bool {

	for col := 0; col < self.file.MaxCol(); col++ {

		data := self.GetValue(row, col, false)

		if data != "" {
			return false
		}
	}

	return true
}

func (self *CSVSheet) GetValue(row, col int, isFloat bool) string {

	rowData := self.file.records[row]

	if col >= len(rowData) {
		return ""
	}

	return rowData[col]
}
