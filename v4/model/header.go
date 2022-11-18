package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
)

type HeaderField struct {
	Col      int
	TypeInfo *DataField // 在类型表中找到对应的类型信息
	tab      *DataTable
}

func (self *HeaderField) String() string {
	return fmt.Sprintf("%s %s %s @%s|%s(%s)", self.TypeInfo.FieldName, self.TypeInfo.FieldType, self.TypeInfo.Comment, self.tab.FileName, self.tab.HeaderType, util.R1C1ToA1(1, self.Col+1))
}

func NewHeaderField(col int) *HeaderField {
	return &HeaderField{
		Col:      col,
		TypeInfo: &DataField{},
	}
}
