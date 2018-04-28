package model

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/table"
)

type SymbolTable struct {
	fields []*table.TableField // 不是具体的类型
}

func (self *SymbolTable) Print() {

	log.Debugln("Symbols:")

	for _, f := range self.fields {
		log.Debugf("%+v", f)
	}

}

func (self *SymbolTable) AddField(tf *table.TableField) {

	if self.FindField(tf.ObjectType, tf.FieldName) != nil {
		panic("Duplicate table field: " + tf.FieldName)
	}

	self.fields = append(self.fields, tf)
}

// 类型是枚举
func (self *SymbolTable) IsEnumKind(objectType string) bool {

	return linq.From(self.EnumNames()).WhereT(func(name string) bool {
		return name == objectType
	}).Count() == 1
}

// 匹配枚举值
func (self *SymbolTable) ResolveEnumValue(objectType, value string) (ret string) {

	linq.From(self.fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType &&
			(tf.Name == value || tf.FieldName == value)
	}).ForEachT(func(types *table.TableField) {

		ret = types.Value

	})

	return
}

// 获取所有的结构体名
func (self *SymbolTable) StructNames() (ret []string) {

	linq.From(self.fields).WhereT(func(tf *table.TableField) bool {

		return tf.Kind == "表头"
	}).SelectT(func(tf *table.TableField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *SymbolTable) EnumNames() (ret []string) {

	linq.From(self.fields).WhereT(func(tf *table.TableField) bool {

		return tf.Kind == "枚举"
	}).SelectT(func(tf *table.TableField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *SymbolTable) Fields(objectType string) (ret []*table.TableField) {

	linq.From(self.fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *SymbolTable) FindField(objectType, name string) (ret *table.TableField) {

	linq.From(self.fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType &&
			(tf.Name == name || tf.FieldName == name)
	}).ForEachT(func(types *table.TableField) {

		ret = types

	})

	return
}

func NewSymbolTable() *SymbolTable {
	return new(SymbolTable)
}
