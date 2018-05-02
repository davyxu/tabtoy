package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/table"
)

type Cell struct {
	Value string
	Row   int // base 0
	Col   int // base 0
	File  string
	Sheet string
}

func (self *Cell) String() string {
	return fmt.Sprintf("'%s' %s|%s(%s)", self.Value, self.File, self.Sheet, util.R1C1ToA1(self.Row+1, self.Col+1))
}

type DataRow []Cell

type DataTable struct {
	HeaderType string // 表名，Index表里定义的类型

	OriginalHeaderType string // HeaderFields对应的ObjectType，KV表为TableField

	FileName string

	SheetName string

	Rows []DataRow

	RawHeader    DataRow
	HeaderFields []*table.TableField // 列索引
}

// 根据列头找到该行对应的值
func (self *DataTable) GetValueByName(row int, name string) (Cell, *table.TableField) {

	for col, tf := range self.HeaderFields {
		if tf.Name == name || tf.FieldName == name {
			return self.Rows[row][col], tf
		}
	}

	return Cell{}, nil
}

// 代码生成专用
func (self *DataTable) GetValue(row, col int) string {

	return self.Rows[row][col].Value
}

// 代码生成专用
func (self *DataTable) GetType(col int) *table.TableField {
	return self.HeaderFields[col]
}

func (self *DataTable) RowCount() int {
	return len(self.Rows)
}

// 添加表头类型
func (self *DataTable) AddHeaderField(types *table.TableField) {
	self.HeaderFields = append(self.HeaderFields, types)
}

func (self *DataTable) HeaderFieldByName(name string) (*table.TableField, int) {

	if name == "" {
		return nil, -1
	}

	for col, f := range self.HeaderFields {
		if f.Name == name || f.FieldName == name {
			return f, col
		}
	}

	return nil, -1
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
