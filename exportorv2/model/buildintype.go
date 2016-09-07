package model

import (
	"github.com/davyxu/tabtoy/proto/tool"
)

type BuildInTypeKind int

const (
	BuildInTypeKind_None BuildInTypeKind = iota
	BuildInTypeKind_Enum
	BuildInTypeKind_Struct
)

type BuildInTypeField struct {
	Name  string
	Value int32
	Meta  tool.FieldMetaV2 // 扩展字段
}

type BuildInType struct {
	Name string
	Kind BuildInTypeKind

	FieldByName   map[string]*BuildInTypeField
	FieldByNumber map[int32]*BuildInTypeField
	Fields        []*BuildInTypeField
}

func (self *BuildInType) Add(def *BuildInTypeField) {

	self.FieldByName[def.Name] = def
	self.FieldByNumber[def.Value] = def
	self.Fields = append(self.Fields, def)
}

func (self *BuildInType) FieldByValueAndMeta(value string) *BuildInTypeField {

	for _, v := range self.FieldByName {

		if v.Name == value {
			return v
		}

		if v.Meta.Alias == value {
			return v
		}

	}

	return nil
}

func NewBuildInType() *BuildInType {
	return &BuildInType{
		FieldByName:   make(map[string]*BuildInTypeField),
		FieldByNumber: make(map[int32]*BuildInTypeField),
	}
}

type BuildInTypeSet struct {
	TypeByName map[string]*BuildInType
}

func (self *BuildInTypeSet) Add(def *BuildInType) {

	self.TypeByName[def.Name] = def
}

func NewBuildInTypeSet() *BuildInTypeSet {
	return &BuildInTypeSet{
		TypeByName: make(map[string]*BuildInType),
	}
}
