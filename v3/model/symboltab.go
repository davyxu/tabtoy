package model

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/table"
	"strconv"
)

type SymbolTable struct {
	typeFields []*table.TypeField // 不是具体的类型
}

func (self *SymbolTable) AddField(tf *table.TypeField) {
	self.typeFields = append(self.typeFields, tf)
}

// 类型是枚举
func (self *SymbolTable) IsEnumKind(tableName, objectType string) bool {

	return linq.From(self.EnumNames()).WhereT(func(name string) bool {
		return name == objectType
	}).Count() == 1
}

// 匹配枚举值
func (self *SymbolTable) ResolveEnumValue(tableName, objectType, value string) (ret string) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.Table == tableName &&
			tf.ObjectType == objectType &&
			(tf.Name == value || tf.FieldName == value)
	}).ForEachT(func(types *table.TypeField) {

		ret = types.DefaultValue

	})

	return
}

// 获取所有的结构体名
func (self *SymbolTable) StructNames() (ret []string) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.DefaultValue == ""
	}).SelectT(func(tf *table.TypeField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *SymbolTable) EnumNames() (ret []string) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.FieldType == "int32" && isNumber(tf.DefaultValue)
	}).SelectT(func(tf *table.TypeField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象在的表名
func (self *SymbolTable) ObjectAtTable(objName string) (ret string) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.ObjectType == objName
	}).SelectT(func(tf *table.TypeField) string {

		return tf.Table
	}).Distinct().ForEachT(func(name string) {

		ret = name

	})

	return
}

// 对象的所有字段
func (self *SymbolTable) Fields(objectType string) (ret []*table.TypeField) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.ObjectType == objectType
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *SymbolTable) QueryType(tableName, headerName string) (ret *table.TypeField) {

	linq.From(self.typeFields).WhereT(func(tf *table.TypeField) bool {

		return tf.Table == tableName &&
			tf.ObjectType == tableName &&
			(tf.Name == headerName || tf.FieldName == headerName)
	}).ForEachT(func(types *table.TypeField) {

		ret = types

	})

	return
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
