package model

import (
	"fmt"
	"github.com/davyxu/golexer"
)

type ObjectFieldType struct {
	ObjectType string
	FieldName  string
	FieldType  string // 去掉[]
	IsArray    bool
	IsStruct   bool

	Meta *golexer.KVPair

	Comment string
}

func (self *ObjectFieldType) String() string {

	return fmt.Sprintf("Object:%s Name: %s Type: %s Array: %v Struct: %v Comment: %s",
		self.ObjectType,
		self.FieldName,
		self.FieldType,
		self.IsArray,
		self.IsStruct,
		self.Comment)

}
