package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
)

type Cell struct {
	Value string
	Row   int // base 0
	Col   int // base 0
	Table *DataTable
}

func (self *Cell) CopyFrom(c *Cell) {
	self.Value = c.Value
	self.Row = c.Row
	self.Col = c.Col
	self.Table = c.Table
}

func (self *Cell) String() string {

	var file, sheet string
	if self.Table != nil {
		file = self.Table.FileName
		sheet = self.Table.SheetName
	}

	return fmt.Sprintf("'%s' @%s|%s(%s)", self.Value, file, sheet, util.R1C1ToA1(self.Row+1, self.Col+1))
}
