package model

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq"
)

type TypeData struct {
	Define *TypeDefine
	Tab    *DataTable // 类型引用的表
	Row    int        // 类型引用的原始数据(DataTable)中的行
}

type TypeTable struct {
	fields []*TypeData
}

func (self *TypeTable) ToJSON(all bool) []byte {

	data, _ := json.MarshalIndent(self.AllFields(all), "", "\t")

	return data
}

func (self *TypeTable) Print(all bool) {

	fmt.Println(string(self.ToJSON(all)))
}

// refData，类型表对应源表的位置信息
func (self *TypeTable) AddField(tf *TypeDefine, data *DataTable, row int) {

	if self.FieldByName(tf.ObjectType, tf.FieldName) != nil {
		panic("Duplicate table field: " + tf.FieldName)
	}

	self.fields = append(self.fields, &TypeData{
		Tab:    data,
		Define: tf,
		Row:    row,
	})
}

func (self *TypeTable) Raw() []*TypeData {
	return self.fields
}

func (self *TypeTable) AllFields(all bool) (ret []*TypeDefine) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		if !all && td.Define.IsBuiltin {
			return false
		}

		return true
	}).SelectT(func(td *TypeData) interface{} {

		return td.Define
	}).ToSlice(&ret)

	return
}

// 类型是枚举
func (self *TypeTable) IsEnumKind(objectType string) bool {

	return linq.From(self.rawEnumNames(true)).WhereT(func(name string) bool {
		return name == objectType
	}).Count() == 1
}

func (self *TypeTable) ResolveEnum(objectType, value string) *TypeData {

	t := self.GetEnumValue(objectType, value)

	if t != nil {
		return t
	}

	enumFields := self.getEnumFields(objectType)

	if len(enumFields) == 0 {
		return nil
	}

	// 默认取第一个
	return enumFields[0]
}

func (self *TypeTable) GetEnumValue(objectType, value string) *TypeData {
	enumFields := self.getEnumFields(objectType)

	if len(enumFields) == 0 {
		return nil
	}

	for _, td := range enumFields {

		if td.Define.Name == value || td.Define.FieldName == value {
			return td
		}

	}

	return nil
}

// 匹配枚举值
func (self *TypeTable) ResolveEnumValue(objectType, value string) string {

	t := self.ResolveEnum(objectType, value)
	if t == nil {
		return ""
	}

	return t.Define.Value
}

func (self *TypeTable) getEnumFields(objectType string) (ret []*TypeData) {

	for _, td := range self.fields {

		if td.Define.ObjectType == objectType {
			ret = append(ret, td)
		}

	}

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

		tf := td.Define

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == TypeUsage_HeaderStruct
	}).SelectT(func(td *TypeData) string {

		return td.Define.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *TypeTable) rawEnumNames(all bool) (ret []string) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		tf := td.Define

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == TypeUsage_Enum
	}).SelectT(func(td *TypeData) string {

		return td.Define.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *TypeTable) AllFieldByName(objectType string) (ret []*TypeDefine) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		return td.Define.ObjectType == objectType
	}).SelectT(func(td *TypeData) *TypeDefine {

		return td.Define
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *TypeTable) FieldByName(objectType, name string) (ret *TypeDefine) {

	linq.From(self.fields).WhereT(func(td *TypeData) bool {

		tf := td.Define

		return tf.ObjectType == objectType &&
			(tf.Name == name || tf.FieldName == name)
	}).ForEachT(func(td *TypeData) {

		ret = td.Define

	})

	return
}

func (self *TypeTable) ObjectExists(objectType string) bool {

	return linq.From(self.fields).WhereT(func(td *TypeData) bool {

		return td.Define.ObjectType == objectType
	}).Count() > 0
}

func NewSymbolTable() *TypeTable {
	return new(TypeTable)
}

var BuiltinSymbolsVisible bool
