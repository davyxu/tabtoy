package tests

import (
	"testing"
)

const (
	TestCSVName = "TestData"
)

// 空类型
func TestEmptyFieldType(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id") // field name
	dataSheet.WriteRow("")   // field type
	dataSheet.WriteRow()     // meta
	dataSheet.WriteRow()     // comment

	emu.MustGotError("TableError.UnknownFieldType 未知字段类型 |  Id   @TestData|TestData(A1)")
}

// 错误的meta
func TestInvalidMeta(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id")    // field name
	dataSheet.WriteRow("int32") // field type
	dataSheet.WriteRow("Drive") // meta
	dataSheet.WriteRow()        // comment

	emu.MustGotError("TableError.UnknownMetaKey 未知的元属性键 | Drive Id int32  @TestData|TestData(A1)")
}

// 空字段名
func TestEmptyFieldName(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id", "", "Id2")           // field name
	dataSheet.WriteRow("int32", "int32", "int32") // field type
	dataSheet.WriteRow()                          // meta
	dataSheet.WriteRow()                          // comment

	emu.MustGotError("TableError.UnknownFieldName 未知字段名 |    @TestData|TestData(B1)")
}

// 表头字段重复
func TestDuplicateHeaderField(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id", "Id")       // field name
	dataSheet.WriteRow("int32", "int32") // field type
	dataSheet.WriteRow()                 // meta
	dataSheet.WriteRow()                 // comment

	emu.MustGotError("TableError.DuplicateHeaderField 表头字段重复 | Id int32  @TestData|TestData(A1)")
}

// 行注释
func TestRowComment(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id", "Str")       // field name
	dataSheet.WriteRow("int32", "string") // field type
	dataSheet.WriteRow()                  // meta
	dataSheet.WriteRow()                  // comment
	dataSheet.WriteRow("#1", "A")
	dataSheet.WriteRow("2", "B")
	dataSheet.WriteRow("#3", "C")
	dataSheet.WriteRow("4", "D")
	dataSheet.WriteRow("#4", "E")

	emu.VerifyData(`{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "",
        	"TestData": [
        		{
        			"Id": 2,
        			"Str": "B"
        		},
        		{
        			"Id": 4,
        			"Str": "D"
        		}
        	]
        }`)
}

// 列注释
func TestColComment(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id", "#Id2", "Id3")       // field name
	dataSheet.WriteRow("int32", "int32", "int32") // field type
	dataSheet.WriteRow()                          // meta
	dataSheet.WriteRow()                          // comment
	dataSheet.WriteRow("1", "2", "3")

	emu.VerifyType(`[
        	{
        		"Usage": 1,
        		"ObjectType": "TestData",
        		"FieldName": "Id",
        		"FieldType": "int32"
        	},
        	{
        		"Usage": 1,
        		"ObjectType": "TestData",
        		"FieldName": "Id3",
        		"FieldType": "int32"
        	}
        ]
`)

	emu.VerifyData(`{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "",
        	"TestData": [
        		{
        			"Id": 1,
        			"Id3": 3
        		}
        	]
        }`)
}

// 数组切割
func TestArraySpliter(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id", "Str")                        // field name
	dataSheet.WriteRow("int32", "string")                  // field type
	dataSheet.WriteRow("ArraySpliter=|", "ArraySpliter=;") // meta
	dataSheet.WriteRow()                                   // comment
	dataSheet.WriteRow("1|2|3", "A;B")
	dataSheet.WriteRow("|4|", "C")

	emu.VerifyData(`{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "",
        	"TestData": [
        		{
        			"Id": [
        				1,
        				2,
        				3
        			],
        			"Str": [
        				"A",
        				"B"
        			]
        		},
        		{
        			"Id": [
        				0,
        				4,
        				0
        			],
        			"Str": [
        				"C"
        			]
        		}
        	]
        }`)
}

// 索引检查
func TestMakeIndex(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Id")        // field name
	dataSheet.WriteRow("int32")     // field type
	dataSheet.WriteRow("MakeIndex") // meta
	dataSheet.WriteRow()            // comment
	dataSheet.WriteRow("1")
	dataSheet.WriteRow("1")
	dataSheet.WriteRow("3")

	emu.MustGotError("TableError.DuplicateValueInMakingIndex 创建索引时发现重复值 | '1' @TestData|TestData(A2)")
}

// 表头类型重复
func TestDuplicateHeaderType(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet1 := emu.CreateDataSheet("data1", "")
	dataSheet1.WriteRow("Id")    // field name
	dataSheet1.WriteRow("int32") // field type
	dataSheet1.WriteRow()        // meta
	dataSheet1.WriteRow()        // comment

	dataSheet2 := emu.CreateDataSheet("data2", "data1")
	dataSheet2.WriteRow("Id")    // field name
	dataSheet2.WriteRow("int32") // field type
	dataSheet2.WriteRow()        // meta
	dataSheet2.WriteRow()        // comment

	emu.MustGotError("TableError.DuplicateHeaderType 表头类型重复 | data1")
}
