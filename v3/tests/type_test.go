package tests

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"testing"
)

// 类型字段重复
func TestDuplicateTypeFieldName(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.MustGotError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type.xlsx|Default(D3)")
}

// 多表中的类型字段重复
func TestDuplicateTypeFieldNameInMultiTypesTable(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type1.xlsx")
	helper.WriteRowValues(indexSheet, "类型表", "", "Type2.xlsx")

	typeSheet1 := emu.CreateDefault("Type1.xlsx")
	helper.WriteTypeTableHeader(typeSheet1)
	helper.WriteRowValues(typeSheet1, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	typeSheet2 := emu.CreateDefault("Type2.xlsx")
	helper.WriteTypeTableHeader(typeSheet2)
	helper.WriteRowValues(typeSheet2, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.MustGotError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type2.xlsx|Default(D2)")
}

// 不填枚举值报错
func TestEnumValueEmpty(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "")

	emu.MustGotError("TableError.EnumValueEmpty 枚举值空 | '' @Type.xlsx|Default(G2)")
}

// 枚举值重复报错
func TestDuplicateEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "1")
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "Arch", "int", "", "1")

	emu.MustGotError("TableError.DuplicateEnumValue 枚举值重复 | '1' @Type.xlsx|Default(G3)")
}

// 枚举值
func TestEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "0")
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "Arch", "int", "", "1")

	helper.WriteRowValues(typeSheet, "表头", "TestData", "角色类型", "Type", "ActorType", "", "")

	dataSheet := emu.CreateDefault("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "角色类型")
	helper.WriteRowValues(dataSheet, "None")
	helper.WriteRowValues(dataSheet, "Arch")

	emu.VerifyData(`
{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "testver",	
	"TestData":[ 
		{ "Type": 0 },
		{ "Type": 1 } 
	]
}
`)
}

//func TestTypeDefineOrder(t *testing.T) {
//
//	mf := NewMemFile()
//	indexSheet := emu.CreateDefault("Index.xlsx")
//
//	helper.WriteIndexTableHeader(indexSheet)
//	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
//
//	typeSheet := emu.CreateDefault("Type.xlsx")
//	helper.WriteTypeTableHeader(typeSheet)
//	helper.WriteRowValues(typeSheet, "表头", "TestData", "角色类型", "Type", "ActorType", "", "")
//	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "0")
//	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "法鸡", "Pharah", "int", "", "1")
//	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "狂鼠", "Junkrat", "int", "", "2")
//
//	if err := VerifyType(mf, ``); err != nil {
//		t.Error(err)
//		t.FailNow()
//	}
//}

func TestBasicType(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "字符串", "String", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "布尔", "Bool", "bool", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "浮点", "Float", "float", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形数组", "IntList", "int", "|", "")

	dataSheet := emu.CreateDefault("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "整形", "字符串", "布尔", "浮点", "整形数组")
	helper.WriteRowValues(dataSheet, "100", "\"hello\"", "true", "3.14159", "1|2|3")

	emu.VerifyGoTypeAndJson(`
{
	"TestData": [
		{
			"Int": 100,
			"String": "\"hello\"",
			"Bool": true,
			"Float": 3.14159,
			"IntList": [
				1,
				2,
				3
			]
		}
	]
}
`)
}

// 禁用索引表和类型表的行
func TestDisableIndexAndTypeRow(t *testing.T) {

	emu := NewTableEmulator(t)

	emu.CreateDefault("TestData1.xlsx")
	emu.CreateDefault("TestData3.xlsx")

	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData1.xlsx")
	helper.WriteRowValues(indexSheet, "#数据表", "", "TestData2.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData3.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData1", "整形", "Int", "int", "", "")
	helper.WriteRowValues(typeSheet, "#表头", "TestData2", "布尔", "Bool", "bool", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData3", "字符串", "String", "string", "", "")

	emu.VerifyType(`
[
	{
		"Kind": 1,
		"ObjectType": "TestData1",
		"Name": "整形",
		"FieldName": "Int",
		"FieldType": "int"
	},
	{

		"Kind": 1,
		"ObjectType": "TestData3",
		"Name": "字符串",
		"FieldName": "String",
		"FieldType": "string"
	}
]
`)
}
