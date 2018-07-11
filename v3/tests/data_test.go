package tests

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"testing"
)

// 行禁用+列禁用
func TestDisableDataRow(t *testing.T) {

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

	dataSheet := emu.CreateDefault("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "整形", "字符串", "#浮点", "布尔")
	helper.WriteRowValues(dataSheet, "100", "\"hello1\"", "1", "")
	helper.WriteRowValues(dataSheet, "200", "\"hello2\"", "2", "true")
	helper.WriteRowValues(dataSheet, "#300", "\"hello3\"", "3", "是")
	helper.WriteRowValues(dataSheet, "400", "\"hello4\"", "4", "")

	emu.VerifyData(`
{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "testver",	
	"TestData":[ 
		{ "Int": 100, "String": "\"hello1\"", "Bool": false, "Float": 0 },
		{ "Int": 200, "String": "\"hello2\"", "Bool": true, "Float": 0 },
		{ "Int": 400, "String": "\"hello4\"", "Bool": false, "Float": 0 } 
	]
}
`)
}

// 表头字段重复
func TestDuplicateHeaderField(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")

	dataSheet := emu.CreateDefault("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "整形", "整形")
	helper.WriteRowValues(dataSheet, "100", "200")

	emu.MustGotError("TableError.DuplicateHeaderField 表头字段重复 | '整形' @TestData.xlsx|Default(A1)")
}

// 数组多列
func TestArrayList(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateDefault("Index.xlsx")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := emu.CreateDefault("Type.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID", "ID", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "技能列表", "SkillList", "int", "|", "")

	dataSheet := emu.CreateDefault("TestData.xlsx")
	helper.WriteRowValues(dataSheet, "ID", "技能列表", "技能列表")
	helper.WriteRowValues(dataSheet, "1", "100", "200")
	helper.WriteRowValues(dataSheet, "2", "", "1") // 多列数组补0

	emu.VerifyData(`
{
			"@Tool": "github.com/davyxu/tabtoy",
			"@Version": "testver",	
			"TestData":[ 
				{ "ID": 1, "SkillList": [100,200] },
				{ "ID": 2, "SkillList": [0,1] } 
			]
		}
`)
}

// TODO KV表测试
