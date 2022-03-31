package model

import (
	"fmt"
	"strings"
)

type HeaderField struct {
	Cell     *Cell       // 表头单元格内容
	TypeInfo *TypeDefine // 在类型表中找到对应的类型信息
}

// String 格式化
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

type HeaderStruct struct {
	Name     string                 // 表头结构体名
	TypeInfo map[string]*TypeDefine // 字段
}

// HadderStructCache 表头结构缓存
var HadderStructCache = make(map[string]*HeaderStruct)

// AddHadderStructCache 添加表头结构数据
func AddHadderStructCache(name string, s *TypeDefine) {
	if _, ok := HadderStructCache[name]; !ok {
		HadderStructCache[name] = &HeaderStruct{
			Name:     name,
			TypeInfo: make(map[string]*TypeDefine),
		}
	}
	// 存放标识名
	HadderStructCache[name].TypeInfo[s.Name] = s
	// 存放字段名
	HadderStructCache[name].TypeInfo[s.FieldName] = s
}
