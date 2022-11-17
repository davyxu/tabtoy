package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"strings"
)

type HeaderField struct {
	Col      int
	TypeInfo *DataField // 在类型表中找到对应的类型信息
	tab      *DataTable
}

func (self *HeaderField) String() string {

	var sb strings.Builder

	fmt.Fprintf(&sb, "Col: %d", self.Col)

	if self.TypeInfo != nil {
		fmt.Fprintf(&sb, "TypeInfo: %+v", self.TypeInfo)
	}

	return sb.String()
}

func (self *HeaderField) Location() string {
	return fmt.Sprintf("@%s|%s(%s)", self.tab.FileName, self.tab.HeaderType, util.R1C1ToA1(0, self.Col))
}

func NewHeaderField(col int, tab *DataTable) *HeaderField {
	return &HeaderField{
		Col:      col,
		TypeInfo: &DataField{},
		tab:      tab,
	}
}
