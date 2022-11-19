package model

import (
	"fmt"
	"strings"
)

// 表格的完整数据，表头有屏蔽时，对应行值为空
type DataTable struct {
	Mode       string
	HeaderType string // Sheet名称, 表名, 结构体名
	FileName   string

	Rows    []*DataRow // 0下标为表头数据
	Headers []*HeaderField
}

// 重复列在表中的索引, 相对于重复列的数量
func (self *DataTable) ArrayFieldCount(field *HeaderField) (ret int) {

	for _, hf := range self.Headers {
		if hf.TypeInfo != nil && hf.TypeInfo.FieldName == field.TypeInfo.FieldName {

			ret++
		}
	}

	return
}

// 模板用，排除表头的数据索引
func (self *DataTable) DataRowIndex() (ret []int) {

	numRows := len(self.Rows)

	if numRows == 0 {
		return
	}

	ret = make([]int, numRows-1)

	// 排除表头数据
	for i := 0; i < numRows-1; i++ {
		ret[i] = i + 1
	}

	return
}

func (self *DataTable) String() string {

	var sb strings.Builder
	sb.WriteString("====DataTable====\n")
	sb.WriteString(fmt.Sprintf("HeaderType: %s\n", self.HeaderType))
	sb.WriteString(fmt.Sprintf("FileName: %s\n", self.FileName))

	// 遍历所有行
	for row, rowData := range self.Rows {

		sb.WriteString(fmt.Sprintf("%d ", row))

		// 遍历一行中的所有列值
		for index, cell := range rowData.Cells() {

			if index > 0 {
				sb.WriteString("/")
			}

			sb.WriteString(cell.Value)

		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func (self *DataTable) AddHeader(header *HeaderField) {
	header.tab = self
	self.Headers = append(self.Headers, header)
}

func (self *DataTable) HeaderByColumn(col int) *HeaderField {

	if col >= len(self.Headers) {
		return nil
	}

	return self.Headers[col]
}

func (self *DataTable) HeaderByName(name string) *HeaderField {
	for _, header := range self.Headers {

		if header.TypeInfo == nil {
			continue
		}

		if header.TypeInfo.FieldName == name {
			return header
		}
	}

	return nil
}

func (self *DataTable) AddRow() (newRow *DataRow) {

	newRow = newDataRow(len(self.Rows), self)

	self.Rows = append(self.Rows, newRow)

	return
}

func (self *DataTable) AddCell(row int) *Cell {

	if row >= len(self.Rows) {
		return nil
	}

	rowData := self.Rows[row]

	return rowData.AddCell()
}

func (self *DataTable) MustGetCell(row, col int) *Cell {

	for len(self.Rows) <= row {
		self.AddRow()
	}

	rowData := self.Rows[row]
	for len(rowData.cells) <= col {
		rowData.AddCell()
	}

	return rowData.Cell(col)
}

// 代码生成专用
func (self *DataTable) GetCell(row, col int) *Cell {

	if row >= len(self.Rows) {
		return nil
	}

	rowData := self.Rows[row]

	if col >= len(rowData.cells) {
		return nil
	}

	return rowData.Cell(col)
}

// 根据列头找到该行对应的值
func (self *DataTable) GetValueByName(row int, name string) *Cell {

	header := self.HeaderByName(name)

	if header == nil {
		return nil
	}

	return self.GetCell(row, header.Col)
}

func NewDataTable() *DataTable {
	return &DataTable{}
}
