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

func (self *TypeTable) ToJSON() []byte {

	data, _ := json.MarshalIndent(self.AllFields(), "", "\t")

	return data
}

func (self *TypeTable) Print() {

	fmt.Println(string(self.ToJSON()))
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

func (self *TypeTable) AllFields() (ret []*TypeDefine) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		if td.Define.IsBuiltin {
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

	for _, tf := range self.fields {
		if tf.Define.Kind == TypeUsage_Enum && objectType == tf.Define.ObjectType {
			return true
		}
	}

	return false
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

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
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

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)

		tf := td.Define

		if !all && tf.IsBuiltin {
			return false
		}

		return tf.Kind == TypeUsage_Enum
	}).Select(func(raw interface{}) interface{} {
		td := raw.(*TypeData)
		return td.Define.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *TypeTable) AllFieldByName(objectType string) (ret []*TypeDefine) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		return td.Define.ObjectType == objectType
	}).Select(func(raw interface{}) interface{} {
		td := raw.(*TypeData)
		return td.Define
	}).ToSlice(&ret)

	return
}

// 数据表中表头对应类型表
func (self *TypeTable) FieldByName(objectType, name string) (ret *TypeDefine) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		tf := td.Define

		return tf.ObjectType == objectType && (tf.Name == name || tf.FieldName == name)
	}).ForEach(func(raw interface{}) {
		td := raw.(*TypeData)
		ret = td.Define
	})

	return
}

func (self *TypeTable) ObjectExists(objectType string) bool {

	return linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		return td.Define.ObjectType == objectType
	}).Count() > 0
}

func NewSymbolTable() *TypeTable {
	return new(TypeTable)
}

var BuiltinSymbolsVisible bool
