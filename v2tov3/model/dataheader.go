package model

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v3/model"
)

type ObjectFieldType struct {
	model.TypeDefine
	Meta *golexer.KVPair
}

func (self *ObjectFieldType) IsArray() bool {
	return self.ArraySplitter != ""
}

//
//import (
//	"fmt"
//	"github.com/davyxu/golexer"
//)
//
//type FieldKind int
//
//const (
//	FieldKind_Primitive FieldKind = 0
//	FieldKind_Enum                = 1
//	FieldKind_Struct              = 2
//)
//
//type ObjectFieldType struct {
//	ObjectType string
//	FieldName  string
//	FieldType  string // 去掉[]
//	Kind       FieldKind
//	IsArray    bool
//
//	Meta *golexer.KVPair
//
//	Comment string
//}
//
//func (self *ObjectFieldType) String() string {
//
//	return fmt.Sprintf("Object:%s Name: %s Type: %s IsArray: %v Kind: %v Comment: %s",
//		self.ObjectType,
//		self.FieldName,
//		self.FieldType,
//		self.IsArray,
//		self.Kind,
//		self.Comment)
//
//}
