package model

import "github.com/davyxu/tabtoy/v3/table"

type DataRow []string

func (self DataRow) Exists(value string) bool {
	for _, v := range self {
		if v == value {
			return true
		}
	}

	return false
}

type DataTable struct {
	Name     string // 表名
	FileName string

	Rows []DataRow

	RawHeader    DataRow
	HeaderFields []*table.TableField // 列索引
}

// 根据列头找到该行对应的值
func (self *DataTable) GetValueByName(row int, name string) (string, *table.TableField) {

	for col, tf := range self.HeaderFields {
		if tf.Name == name || tf.FieldName == name {
			return self.Rows[row][col], tf
		}
	}

	return "", nil
}

// 代码生成专用
func (self *DataTable) GetValue(row, col int) string {

	return self.Rows[row][col]
}

// 代码生成专用
func (self *DataTable) GetType(col int) *table.TableField {
	return self.HeaderFields[col]
}

func (self *DataTable) HeaderFieldCount() int {
	return len(self.RawHeader)
}

func (self *DataTable) RowCount() int {
	return len(self.Rows)
}

// 添加表头类型
func (self *DataTable) AddHeaderField(types *table.TableField) {
	self.HeaderFields = append(self.HeaderFields, types)
}

// 添加行数据
func (self *DataTable) AddRow(row DataRow) {

	self.Rows = append(self.Rows, row)
}

// 获取一整行数据
func (self *DataTable) GetDataRow(row int) DataRow {
	return self.Rows[row]
}

func NewDataTable() *DataTable {
	return &DataTable{}
}
