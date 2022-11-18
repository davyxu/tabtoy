package model

import (
	"encoding/json"
	"fmt"
	"github.com/ahmetb/go-linq"
)

type TypeData struct {
	Define *DataField
	Tab    *DataTable // 类型引用的表
	Row    int        // 类型引用的原始数据(DataTable)中的行
}

type TypeManager struct {
	fields []*TypeData
}

func NewTypeManager() *TypeManager {
	return new(TypeManager)
}

func (self *TypeManager) ToJSON() []byte {

	data, _ := json.MarshalIndent(self.AllFields(), "", "\t")

	return data
}

func (self *TypeManager) Print() {

	fmt.Println(string(self.ToJSON()))
}

// refData，类型表对应源表的位置信息
func (self *TypeManager) AddField(tf *DataField, data *DataTable, refRow int) {

	if self.FieldByName(tf.ObjectType, tf.FieldName) != nil {
		panic("Duplicate table field: " + tf.FieldName)
	}

	self.fields = append(self.fields, &TypeData{
		Tab:    data,
		Define: tf,
		Row:    refRow,
	})
}

func (self *TypeManager) Raw() []*TypeData {
	return self.fields
}

func (self *TypeManager) AllFields() (ret []*DataField) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		return true
	}).SelectT(func(td *TypeData) interface{} {

		return td.Define
	}).ToSlice(&ret)

	return
}

// 类型是枚举
func (self *TypeManager) IsEnumKind(objectType string) bool {

	for _, tf := range self.fields {
		if tf.Define.Usage == FieldUsage_Enum && objectType == tf.Define.ObjectType {
			return true
		}
	}

	return false
}

func (self *TypeManager) ResolveEnum(objectType, value string) *TypeData {

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

func (self *TypeManager) GetEnumValue(objectType, value string) *TypeData {
	enumFields := self.getEnumFields(objectType)

	if len(enumFields) == 0 {
		return nil
	}

	for _, td := range enumFields {

		if td.Define.FieldName == value {
			return td
		}

	}

	return nil
}

// 匹配枚举值
func (self *TypeManager) ResolveEnumValue(objectType, value string) string {

	t := self.ResolveEnum(objectType, value)
	if t == nil {
		return ""
	}

	return t.Define.Value
}

func (self *TypeManager) getEnumFields(objectType string) (ret []*TypeData) {

	for _, td := range self.fields {

		if td.Define.ObjectType == objectType {
			ret = append(ret, td)
		}

	}

	return
}

func (self *TypeManager) EnumNames() (ret []string) {

	return self.rawEnumNames(true)
}

func (self *TypeManager) StructNames() (ret []string) {

	return self.rawStructNames(true)
}

// 获取所有的结构体名
func (self *TypeManager) rawStructNames(all bool) (ret []string) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		tf := td.Define

		if !all {
			return false
		}

		return tf.Usage == FieldUsage_Struct
	}).SelectT(func(td *TypeData) string {

		return td.Define.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 获取所有的枚举名
func (self *TypeManager) rawEnumNames(all bool) (ret []string) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)

		tf := td.Define

		if !all {
			return false
		}

		return tf.Usage == FieldUsage_Enum
	}).Select(func(raw interface{}) interface{} {
		td := raw.(*TypeData)
		return td.Define.ObjectType
	}).Distinct().ToSlice(&ret)

	return
}

// 对象的所有字段
func (self *TypeManager) AllFieldByName(objectType string) (ret []*DataField) {

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
func (self *TypeManager) FieldByName(objectType, name string) (ret *DataField) {

	linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		tf := td.Define

		return tf.ObjectType == objectType && tf.FieldName == name
	}).ForEach(func(raw interface{}) {
		td := raw.(*TypeData)
		ret = td.Define
	})

	return
}

func (self *TypeManager) ObjectExists(objectType string) bool {

	return linq.From(self.fields).Where(func(raw interface{}) bool {
		td := raw.(*TypeData)
		return td.Define.ObjectType == objectType
	}).Count() > 0
}
