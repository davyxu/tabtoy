package model

import "sort"

type GlobalChecker interface {
	CheckValueRepeat(fd *FieldDescriptor, value string) bool

	GlobalFileDesc() *FileDescriptor
}

type FieldValue struct {
	FieldDef           *FieldDescriptor
	RawValue           string
	R                  int
	C                  int
	SheetName          string
	FileName           string
	FieldRepeatedCount int // repeated拆成多列时, 这样重复列的数量
}

// 对应record
type LineData struct {
	Values []*FieldValue
}

func (self *LineData) Len() int {
	return len(self.Values)
}

func (self *LineData) Swap(i, j int) {
	self.Values[i], self.Values[j] = self.Values[j], self.Values[i]
}

func (self *LineData) Less(i, j int) bool {

	a := self.Values[i]
	b := self.Values[j]

	// repeated字段分多个单元格导出时, 由于进行数据排序, 所以需要增加对列排序因子保证最终数据正确性
	return a.FieldDef.Order+int32(a.C) < b.FieldDef.Order+int32(b.C)
}

func (self *LineData) Add(fv *FieldValue) {

	self.Values = append(self.Values, fv)
}

func NewLineData() *LineData {
	return new(LineData)
}

// 对应table
type DataModel struct {
	Lines []*LineData
}

func (self *DataModel) Add(data *LineData) {

	sort.Sort(data)

	self.Lines = append(self.Lines, data)
}

func NewDataModel() *DataModel {
	return new(DataModel)
}
