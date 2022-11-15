package util

import (
	"bytes"
	"encoding/csv"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
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

func ConvUTF8ToGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func (self *CSVFile) Transform(convf func(s []byte) ([]byte, error)) {

	for row, rowData := range self.records {

		for col, cell := range rowData {

			if cell != "" {
				data, _ := convf([]byte(cell))
				self.records[row][col] = string(data)
			}
		}

	}
}

func (self *CSVFile) Save(filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	err = writer.WriteAll(self.records)

	if err != nil {
		return err
	}

	writer.Flush()

	return nil
}

func (self *CSVFile) Load(fileName string) error {

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	self.Name = strings.TrimSuffix(fileName, filepath.Ext(fileName))

	self.records, err = csv.NewReader(bytes.NewReader(data)).ReadAll()

	if err != nil {
		return err
	}

	// 自动探测非utf-8编码, 转换
	_, codecName, _ := charset.DetermineEncoding(data, "")
	if codecName != "utf-8" {
		self.Transform(ConvGBKToUTF8)
	}

	return nil
}

func NewCSVFile() *CSVFile {

	self := &CSVFile{}

	// 默认创建sheet
	self.sheet = &CSVSheet{file: self}

	return self
}

type CSVSheet struct {
	file *CSVFile
}

func (self *CSVSheet) Name() string {
	return self.file.Name
}

func (self *CSVSheet) MaxColumn() int {
	return self.file.MaxCol()
}

func (self *CSVSheet) IsRowEmpty(row, maxCol int) bool {

	if maxCol == -1 {
		maxCol = self.file.MaxCol()
	}

	for col := 0; col < maxCol; col++ {

		data := self.GetValue(row, col, nil)

		if data != "" {
			return false
		}
	}

	return true
}

func (self *CSVSheet) GetValue(row, col int, opt *ValueOption) string {

	if row >= len(self.file.records) {
		return ""
	}

	rowData := self.file.records[row]

	if col >= len(rowData) {
		return ""
	}

	return rowData[col]
}

func (self *CSVSheet) SetValue(row, col int, value string) bool {
	if row >= len(self.file.records) {
		return false
	}
	rowData := self.file.records[row]

	if col >= len(rowData) {
		return false
	}

	rowData[col] = value
	return true
}

func (self *CSVSheet) WriteRow(valueList ...string) {

	var rowData []string
	for _, str := range valueList {

		rowData = append(rowData, str)
	}

	self.file.records = append(self.file.records, rowData)
}
