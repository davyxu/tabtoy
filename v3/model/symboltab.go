package model

import (
	"github.com/ahmetb/go-linq"
	"strconv"
)

type SymbolTable struct {
	Types []*ObjectTypes
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func (self *SymbolTable) IsEnumKind(tableName, objectType string) bool {

	return linq.From(self.Types).WhereT(func(types *ObjectTypes) bool {

		return types.Table == tableName &&
			types.ObjectType == objectType &&
			(types.FieldType == "int32" && isNumber(types.DefaultValue)) // 字段类型都为int32，且默认值(枚举值)都
	}).Count() > 0 // 字段数量
}

func (self *SymbolTable) GetEnumValue(tableName, objectType, value string) (ret string) {

	linq.From(self.Types).WhereT(func(types *ObjectTypes) bool {

		return types.Table == tableName &&
			types.ObjectType == objectType &&
			(types.Name == value || types.FieldName == value)
	}).ForEachT(func(types *ObjectTypes) {

		ret = types.DefaultValue

	})

	return
}

func (self *SymbolTable) QueryType(tableName, headerName string) (ret *ObjectTypes) {

	linq.From(self.Types).WhereT(func(types *ObjectTypes) bool {

		return types.Table == tableName &&
			types.ObjectType == tableName &&
			(types.Name == headerName || types.FieldName == headerName)
	}).ForEachT(func(types *ObjectTypes) {

		ret = types

	})

	return
}
