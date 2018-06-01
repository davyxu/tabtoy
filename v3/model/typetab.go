package model

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/table"
)

type TypeData struct {
	Type *table.TableField
	Tab  *DataTable // 类型引用的表
	Row  int        // 类型引用的原始数据(DataTable)中的行
}

type TypeTable struct {
	fields []*TypeData
}

func (self *TypeTable) ToJSON() []byte {

	data, _ := json.MarshalIndent(self.AllFields(true), "", "\t")

	return data
}

func (self *TypeTable) Print() {

	fmt.Println(string(self.ToJSON()))
}

// refData，类型表对应源表的位置信息
func (self *TypeTable) AddField(tf *table.TableField, data *DataTable, row int) {

	if self.FieldByName(tf.ObjectType, tf.FieldName) != nil {
		panic("Duplicate table field: " + tf.FieldName)
	}

	self.fields = append(self.fields, &TypeData{
		Tab:  data,
		Type: tf,
		Row:  row,
	})
}

func (self *TypeTable) Raw() []*TypeData {
	return self.fields
}

func (self *TypeTable) AllFields(all bool) (ret []*table.TableField) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		if !all && td.Type.IsBuiltin {
			return false
		}

		return true
	}).SelectT(func(td *TypeData) interface{} {

		return td.Type
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

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		return td.Type.ObjectType == objectType &&
			(td.Type.Name == value || td.Type.FieldName == value)
	}).ForEachT(func(td *TypeData) {

		ret = td.Type.Value

	})

	return
}

func (self *TypeTable) EnumNames() (ret []string) {

	return self.rawEnumNames(BuiltinSymbolsVisible)
}

func (self *TypeTable) StructNames() (ret []string) {

	return self.rawStructNames(BuiltinSymbolsVisible)
}

// 获取所有的结构体名
func (self *TypeTable) rawStructNames(all bool) (ret []string) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		tf := td.Type

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == table.TableKind_HeaderStruct
	}).SelectT(func(td *TypeData) string {

		return td.Type.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *TypeTable) rawEnumNames(all bool) (ret []string) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		tf := td.Type

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == table.TableKind_Enum
	}).SelectT(func(td *TypeData) string {

		return td.Type.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *TypeTable) AllFieldByName(objectType string) (ret []*table.TableField) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		return td.Type.ObjectType == objectType
	}).SelectT(func(td *TypeData) *table.TableField {

		return td.Type
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *TypeTable) FieldByName(objectType, name string) (ret *table.TableField) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		tf := td.Type

		return tf.ObjectType == objectType &&
			(tf.Name == name || tf.FieldName == name)
	}).ForEachT(func(td *TypeData) {

		ret = td.Type

	})

	return
}

func (self *TypeTable) ObjectExists(objectType string) bool {

	return linq.From(self.fields).WhereT(func(td *TypeData) bool {

		return td.Type.ObjectType == objectType
	}).Count() > 0
}

func NewSymbolTable() *TypeTable {
	return new(TypeTable)
}

var BuiltinSymbolsVisible bool
