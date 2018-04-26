package model

import "github.com/davyxu/tabtoy/v3/table"

type DataRow []string

type DataTable struct {
	name   string             // 表名
	header []*table.TypeField // 列索引
	rows   []DataRow
}

func (self *DataTable) Header() []*table.TypeField {
	return self.header
}

func (self *DataTable) Rows() []DataRow {
	return self.rows
}

func (self *DataTable) Name() string {
	return self.name
}

func (self *DataTable) MaxColumns() int {
	return len(self.header)
}

func (self *DataTable) AddHeader(types *table.TypeField) {
	self.header = append(self.header, types)
}

// row 从0开始，相对DataTable的索引
func (self *DataTable) AddRow(row, col int, value string) {

	var rowData DataRow
	if row < len(self.rows) {
		rowData = self.rows[row]
	} else {
		rowData = make(DataRow, len(self.header))
		self.rows = append(self.rows, rowData)
	}

	rowData[col] = value
}

func (self *DataTable) GetValue(row, col int) string {

	return self.rows[row][col]
}

func (self *DataTable) GetType(col int) *table.TypeField {
	return self.header[col]
}

func NewDataTable(name string) *DataTable {
	return &DataTable{
		name: name,
	}
}
