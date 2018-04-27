package model

import "github.com/davyxu/tabtoy/v3/table"

type DataRow []string

type DataTable struct {
	name      string // 表名
	rawHeader DataRow
	rows      []DataRow

	headerField []*table.TypeField // 列索引
}

func (self *DataTable) Header() []*table.TypeField {
	return self.headerField
}

func (self *DataTable) Rows() []DataRow {
	return self.rows
}

func (self *DataTable) Name() string {
	return self.name
}

func (self *DataTable) RawHeader() DataRow {
	return self.rawHeader
}

func (self *DataTable) HeaderFieldCount() int {
	return len(self.rawHeader)
}

func (self *DataTable) RowCount() int {
	return len(self.rows)
}

func (self *DataTable) AddHeaderField(types *table.TypeField) {
	self.headerField = append(self.headerField, types)
}

func (self *DataTable) AddRow(row DataRow) {

	self.rows = append(self.rows, row)
}

func (self *DataTable) GetDataRow(row int) DataRow {
	return self.rows[row]
}

func (self *DataTable) GetValue(row, col int) string {

	return self.rows[row][col]
}

func (self *DataTable) GetType(col int) *table.TypeField {
	return self.headerField[col]
}

func NewDataTable(name string, rawheader DataRow) *DataTable {
	return &DataTable{
		name:      name,
		rawHeader: rawheader,
	}
}
