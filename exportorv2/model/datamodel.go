package model

import "sort"

type FieldValue struct {
	FieldDef  *FieldDescriptor
	RawValue  string
	R         int
	C         int
	File      interface{}
	SheetName string
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

	return a.FieldDef.Order < b.FieldDef.Order
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
