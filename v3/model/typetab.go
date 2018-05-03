package model

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/table"
)

type TypeTable struct {
	Fields []*table.TableField // 不是具体的类型
}

func (self *TypeTable) ToJSON() []byte {

	data, _ := json.MarshalIndent(&struct {
		Fields []*table.TableField
	}{
		Fields: self.UserFields(),
	}, "", "\t")

	return data
}

func (self *TypeTable) Print() {

	fmt.Println(string(self.ToJSON()))
}

func (self *TypeTable) AddField(tf *table.TableField) {

	if self.FieldByName(tf.ObjectType, tf.FieldName) != nil {
		panic("Duplicate table field: " + tf.FieldName)
	}

	self.Fields = append(self.Fields, tf)
}

//
func (self *TypeTable) UserFields() (ret []*table.TableField) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {
		return !tf.IsBuiltin
	}).ToSlice(&ret)

	return
}

// 类型是枚举
func (self *TypeTable) IsEnumKind(objectType string) bool {

	return linq.From(self.rawEnumNames(true)).WhereT(func(name string) bool {
		return name == objectType
	}).Count() == 1
}

// 匹配枚举值
func (self *TypeTable) ResolveEnumValue(objectType, value string) (ret string) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType &&
			(tf.Name == value || tf.FieldName == value)
	}).ForEachT(func(types *table.TableField) {

		ret = types.Value

	})

	return
}

func (self *TypeTable) EnumNames() (ret []string) {

	return self.rawEnumNames(UseAllBuiltinSymbols)
}

func (self *TypeTable) StructNames() (ret []string) {

	return self.rawStructNames(UseAllBuiltinSymbols)
}

// 获取所有的结构体名
func (self *TypeTable) rawStructNames(all bool) (ret []string) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == table.TableKind_HeaderStruct
	}).SelectT(func(tf *table.TableField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *TypeTable) rawEnumNames(all bool) (ret []string) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == table.TableKind_Enum
	}).SelectT(func(tf *table.TableField) string {

		return tf.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *TypeTable) AllFieldByName(objectType string) (ret []*table.TableField) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *TypeTable) FieldByName(objectType, name string) (ret *table.TableField) {

	linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType &&
			(tf.Name == name || tf.FieldName == name)
	}).ForEachT(func(types *table.TableField) {

		ret = types

	})

	return
}

func (self *TypeTable) ObjectExists(objectType string) bool {

	return linq.From(self.Fields).WhereT(func(tf *table.TableField) bool {

		return tf.ObjectType == objectType
	}).Count() > 0
}

func NewSymbolTable() *TypeTable {
	return new(TypeTable)
}

var UseAllBuiltinSymbols bool
