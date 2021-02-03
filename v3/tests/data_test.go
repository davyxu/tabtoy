package tests

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"testing"
)

// 行禁用+列禁用
func TestDisableDataRow(t *testing.T) {

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

	dataSheet := emu.CreateCSVFile("TestData")
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
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "整形", "整形")
	helper.WriteRowValues(dataSheet, "100", "200")

	emu.MustGotError("TableError.DuplicateHeaderField 表头字段重复 | '整形' @TestData|(A1)")
}

// 数组多列
func TestArrayList(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "TestData", "TestData")
	helper.WriteRowValues(indexSheet, "数据表", "TestData", "TestData2")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID", "ID", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "技能列表", "SkillList", "int", "|", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "名字列表", "NameList", "string", "|", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID列表", "IDList", "int32", "|", "") // 单列

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "ID", "技能列表", "技能列表", "技能列表", "名字列表", "名字列表", "ID列表")
	helper.WriteRowValues(dataSheet, "1", "100", "200", "300", "", "", "1|2")
	helper.WriteRowValues(dataSheet, "2", "1", "", "3", "", "", "")   // 多列数组补0
	helper.WriteRowValues(dataSheet, "3", "", "20", "30", "", "", "") // 多列数组补0
	helper.WriteRowValues(dataSheet, "4", "", "", "", "", "", "")     // 多列数组补0

	// 数组列多列情况, 又在不同的表中可能存在不同数量列的情况
	dataSheet2 := emu.CreateCSVFile("TestData2")
	helper.WriteRowValues(dataSheet2, "ID", "技能列表", "技能列表", "技能列表")
	helper.WriteRowValues(dataSheet2, "5", "1", "1")
	helper.WriteRowValues(dataSheet2, "6", "2", "1")

	emu.VerifyData(`
{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "testver",	
        	"TestData":[ 
        		{ "ID": 1, "SkillList": [100,200,300], "NameList": ["",""], "IDList": [1, 2] },
        		{ "ID": 2, "SkillList": [1,0,3], "NameList": ["",""], "IDList": [] },
        		{ "ID": 3, "SkillList": [0,20,30], "NameList": ["",""], "IDList": [] },
        		{ "ID": 4, "SkillList": [0,0,0], "NameList": ["",""], "IDList": [] },
        		{ "ID": 5, "SkillList": [1,1,0], "NameList": [], "IDList": [] },
        		{ "ID": 6, "SkillList": [2,1,0], "NameList": [], "IDList": [] } 
        	]
        }
`)
}

func TestArrayMultiColumnDefineNotMatch(t *testing.T) {
	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "TestData", "TestData")
	helper.WriteRowValues(indexSheet, "数据表", "TestData", "TestData2")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID", "ID", "int", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TestData", "技能列表", "SkillList", "int", "|", "")

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "ID", "技能列表", "技能列表")
	helper.WriteRowValues(dataSheet, "1", "100", "200")

	// 数组列多列情况, 又在不同的表中可能存在不同数量列的情况
	dataSheet2 := emu.CreateCSVFile("TestData2")
	helper.WriteRowValues(dataSheet2, "ID", "技能列表")
	helper.WriteRowValues(dataSheet2, "5", "")
	helper.WriteRowValues(dataSheet2, "6", "")

	emu.MustGotError("TableError.ArrayMultiColumnDefineNotMatch 数组类型多列跨表定义不一致 | ")
}

// 重复性检查
func TestRepeatCheck(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID", "ID", "int", "", "", "true")

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "ID")
	helper.WriteRowValues(dataSheet, "1")
	helper.WriteRowValues(dataSheet, "1") // 多列数组补0

	emu.MustGotError("TableError.DuplicateValueInMakingIndex 创建索引时发现重复值 | '1' @TestData|(A3)")
}

// TODO KV表测试

// 单元格中有切割符时, 重复列的拆分符需要单独设置
func TestArraySpliter(t *testing.T) {

	emu := NewTableEmulator(t)
	indexSheet := emu.CreateCSVFile("Index")

	helper.WriteIndexTableHeader(indexSheet)
	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")

	typeSheet := emu.CreateCSVFile("Type")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "TestData", "Week", "Week", "string", "$", "", "true")

	dataSheet := emu.CreateCSVFile("TestData")
	helper.WriteRowValues(dataSheet, "Week", "Week")
	helper.WriteRowValues(dataSheet, "1|2|3", "4|5|6")

	emu.VerifyData(`
{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "testver",	
        	"TestData":[ 
        		{ "Week": ["1|2|3", "4|5|6"] } 
        	]
        }
`)
}

//// 单元格类型与期望类型不匹配时
//func TestMissMatchingType(t *testing.T) {
//
//	emu := NewTableEmulator(t)
//	indexSheet := emu.CreateCSVFile("Index")
//
//	helper.WriteIndexTableHeader(indexSheet)
//	helper.WriteRowValues(indexSheet, "类型表", "", "Type")
//	helper.WriteRowValues(indexSheet, "数据表", "", "TestData")
//
//	typeSheet := emu.CreateCSVFile("Type")
//	helper.WriteTypeTableHeader(typeSheet)
//	helper.WriteRowValues(typeSheet, "表头", "TestData", "ID", "ID", "int", "", "", "true")
//
//	dataSheet := emu.CreateCSVFile("TestData")
//	helper.WriteRowValues(dataSheet, "ID")
//	helper.WriteRowValues(dataSheet, "中文")
//
//	emu.VerifyData(`
//{
//        	"@Tool": "github.com/davyxu/tabtoy",
//        	"@Version": "testver",
//        	"TestData":[
//        		{ "Week": ["1|2|3", "4|5|6"] }
//        	]
//        }
//`)
//}
