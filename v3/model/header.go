package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/table"
	"strings"
)

type HeaderField struct {
	Cell     *Cell             // 表头单元格内容
	TypeInfo *table.TableField // 在类型表中找到对应的类型信息
}

func (self *HeaderField) String() string {

	var sb strings.Builder

	if self.Cell != nil {
		sb.WriteString("Cell: ")
		sb.WriteString(self.Cell.String())
	}

	if self.TypeInfo != nil {
		sb.WriteString("TypeInfo: ")
		sb.WriteString(fmt.Sprintf("%+v", self.TypeInfo))
	}

	return sb.String()
}
