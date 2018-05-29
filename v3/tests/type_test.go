package tests

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"testing"
)

// 类型字段重复
func TestDuplicateTypeFieldName(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.Create("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.Create("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.VerifyError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type.xlsx|Default(D3)")
}

// 多表中的类型字段重复
func TestDuplicateTypeFieldNameInMultiTypesTable(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.Create("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type1.xlsx")
	helper.WriteRowValues(indexSheet, "类型表", "", "Type2.xlsx")

	typeSheet1 := emu.Create("Type1.xlsx")
	helper.WriteTypeTableHeader(typeSheet1)
	helper.WriteRowValues(typeSheet1, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	typeSheet2 := emu.Create("Type2.xlsx")
	helper.WriteTypeTableHeader(typeSheet2)
	helper.WriteRowValues(typeSheet2, "表头", "SumeHead", "某种类型", "None", "int", "", "")

	emu.VerifyError("TableError.DuplicateTypeFieldName 类型表字段重复 | 'None' @Type2.xlsx|Default(D2)")
}

// 不填枚举值报错
func TestEnumValueEmpty(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.Create("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.Create("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "")

	emu.VerifyError("TableError.EnumValueEmpty 枚举值空 | '' @Type.xlsx|Default(G2)")
}

// 枚举值重复报错
func TestDuplicateEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.Create("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := emu.Create("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "None", "int", "", "1")
	helper.WriteRowValues(typeSheet, "枚举", "ActorType", "", "Arch", "int", "", "1")

	emu.VerifyError("TableError.DuplicateEnumValue 枚举值重复 | '1' @Type.xlsx|Default(G3)")
}

//func TestTypeDefineOrder(t *testing.T) {
//
//	mf := NewMemFile()
//	indexSheet := emu.Create("Index.xlsx")
//
//	helper.WriteIndexTableHeader(indexSheet)
//	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
//
//	typeSheet := emu.Create("Type.xlsx")
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

func TestComplete(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.Create("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := emu.Create("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "字符串", "String", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "布尔", "Bool", "bool", "", "")

	dataSheet := emu.Create("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "整形", "字符串", "布尔")
	helper.WriteRowValues(dataSheet, "100", "\"hello\"", "true")

	if err := emu.VerifyLauncherJson(`
{
	"TestData": [
		{
			"Int": 100,
			"String": "\"hello\"",
			"Bool": true
		}
	]
}

`); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
