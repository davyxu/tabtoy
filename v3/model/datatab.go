package model

import "github.com/davyxu/tabtoy/v3/table"

type DataRow []string

type DataTable struct {
	name      string // 表名
	rawHeader DataRow
	rows      []DataRow

	headerField []*table.TableField // 列索引
}

// 代码生成专用
func (self *DataTable) Header() []*table.TableField {
	return self.headerField
}

// 代码生成专用
func (self *DataTable) Rows() []DataRow {
	return self.rows
}

// 根据列头找到该行对应的值
func (self *DataTable) GetValueByName(row int, name string) (string, *table.TableField) {

	for col, tf := range self.headerField {
		if tf.Name == name || tf.FieldName == name {
			return self.rows[row][col], tf
		}
	}

	return "", nil
}

// 代码生成专用
func (self *DataTable) GetValue(row, col int) string {

	return self.rows[row][col]
}

// 代码生成专用
func (self *DataTable) GetType(col int) *table.TableField {
	return self.headerField[col]
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

// 添加表头类型
func (self *DataTable) AddHeaderField(types *table.TableField) {
	self.headerField = append(self.headerField, types)
}

// 添加行数据
func (self *DataTable) AddRow(row DataRow) {

	self.rows = append(self.rows, row)
}

// 获取一整行数据
func (self *DataTable) GetDataRow(row int) DataRow {
	return self.rows[row]
}

func NewDataTable(name string, rawheader DataRow) *DataTable {
	return &DataTable{
		name:      name,
		rawHeader: rawheader,
	}
}
