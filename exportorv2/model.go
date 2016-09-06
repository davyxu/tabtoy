package exportorv2

import (
	"fmt"
)

type cellCoord struct {
	col int
	row int
}

// 一个单元格
type CellContext struct {
	*FieldDefine

	Value string // 转换输出的值

	ValueList []string

	rawByCellCoord map[cellCoord]string // 引用的单元格, 值为原始值
}

// base 0
func (self *CellContext) addRefCell(row, col int, raw string) {
	self.rawByCellCoord[cellCoord{col: col, row: row}] = raw
}

// base 0
func (self *CellContext) findRefCell(row, col int) (string, bool) {
	if v, ok := self.rawByCellCoord[cellCoord{col: col, row: row}]; ok {
		return v, true
	}

	return "", false
}

func (self CellContext) String() string {

	var buildInType string
	if self.BuildInType != nil {
		buildInType = fmt.Sprintf("(%s)", self.BuildInType.Name)
	}

	var valuePresent string
	if self.IsRepeated {
		valuePresent = fmt.Sprintf("valuelist: %v", self.ValueList)
	} else {
		valuePresent = fmt.Sprintf("value: %v", self.Value)
	}

	return fmt.Sprintf("name: %s type: %s%s %s", self.Name, FieldTypeToString(self.Type), buildInType, valuePresent)
}

func newCellContext(def *FieldDefine) *CellContext {
	return &CellContext{
		FieldDefine:    def,
		rawByCellCoord: make(map[cellCoord]string),
	}
}

type Record struct {
	cellByFD map[*FieldDefine]*CellContext
	cells    []*CellContext
}

func (self *Record) NewContextByDefine(def *FieldDefine) *CellContext {

	// 如果这个单元格数据有, 使用已经有的定义, 因为字段不会重复
	// 主要处理repeated散开的case
	if exist, ok := self.cellByFD[def]; ok {
		return exist
	}

	ctx := newCellContext(def)

	self.cellByFD[def] = ctx
	self.cells = append(self.cells, ctx)

	return ctx
}

func newRecord() *Record {
	return &Record{
		cellByFD: make(map[*FieldDefine]*CellContext),
	}
}
