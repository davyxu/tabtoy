package tests

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"testing"
)

// 类型字段重复
func TestDuplicateTypeFieldName(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.MustGotError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type|(D3) SumeHead None")
}

// 多表中的类型字段重复
func TestDuplicateTypeFieldNameInMultiTypesTable(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type1")
	helper.WriteRowValues(indexSheet, "类型表", "", "Type2")

	typeSheet1 := emu.CreateCSVFile("Type1")
	helper.WriteTypeTableHeader(typeSheet1)
	helper.WriteRowValues(typeSheet1, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	typeSheet2 := emu.CreateCSVFile("Type2")
	helper.WriteTypeTableHeader(typeSheet2)
	helper.WriteRowValues(typeSheet2, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.MustGotError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type2|(D2) SumeHead None")
}

// 不填枚举值报错
func TestEnumValueEmpty(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "")

	emu.MustGotError("TableError.EnumValueEmpty 枚举值空 | '' @Type|(G2)")
}

// 枚举值重复报错
func TestDuplicateEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "1")
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "Arch", "int", "", "1")

	emu.MustGotError("TableError.DuplicateEnumValue 枚举值重复 | '1' @Type|(G3)")
}

// 枚举值
func TestEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "0")
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "Arch", "int", "", "1")

	helper.WriteRowValues(typeSheet, "表头", "TestData", "角色类型", "Type", "ActorType", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "索引", "Index", "int", "", "")

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "索引", "角色类型")
	helper.WriteRowValues(dataSheet, "", "Arch")
	helper.WriteRowValues(dataSheet, "", "None")
	helper.WriteRowValues(dataSheet, "3", "")

	emu.VerifyData(`
{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "testver",	
	"TestData":[ 
		{ "Type": 1, "Index": 0 },
		{ "Type": 0, "Index": 0 },
		{ "Type": 0, "Index": 3 }
	]
}
`)
}

func TestBasicType(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "字符串", "String", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "布尔", "Bool", "bool", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "浮点", "Float", "float", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形数组", "IntList", "int", "|", "")

	dataSheet := emu.CreateCSVFile("TestData")
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

	emu.CreateCSVFile("TestData1")
	emu.CreateCSVFile("TestData3")

	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData1")
	helper.WriteRowValues(indexSheet, "#数据表", "", "TestData2")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData3")

	typeSheet := emu.CreateCSVFile("Type")
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
