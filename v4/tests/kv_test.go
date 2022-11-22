package tests

import "testing"

// KV表空类型
func TestKVEmptyFieldType(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("Id", "", "", "", "")

	emu.MustGotError("TableError.UnknownFieldType 未知字段类型 |  '' @TestData|(B2)")
}

// KV表错误的meta
func TestKVInvalidMeta(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("Id", "int32", "", "", "Drive")

	emu.MustGotError("TableError.UnknownMetaKey 未知的元属性键 | Drive 'Drive' @TestData|(E2)")
}

// KV表空字段名
func TestKVEmptyFieldName(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("", "int32", "", "", "")

	emu.MustGotError("TableError.UnknownFieldName 未知字段名 | '' @TestData|(A2)")
}

// KV表字段重复
func TestKVDuplicatFieldName(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("Id", "int32", "", "", "")
	dataSheet.WriteRow("Id", "int32", "", "", "")

	emu.MustGotError("TableError.DuplicateKVField KV表字段重复 | 'Id' @TestData|(A3)")
}

// KV表行注释
func TestKVRowComment(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("Id", "int32", "1", "", "")
	dataSheet.WriteRow("#Id2", "int32", "2", "", "")
	dataSheet.WriteRow("Id3", "int32", "3", "", "")
	dataSheet.WriteRow("#Id4", "int32", "4", "", "")

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
        ]`)

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

// KV表数组切割
func TestKVArraySpliter(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet := emu.CreateKVSheet(TestCSVName, "")
	// Key	Type	Value	Comment	Meta
	dataSheet.WriteRow("Id", "int32", "1|2|3", "", "ArraySpliter=|")

	emu.VerifyData(`{
        	"@Tool": "github.com/davyxu/tabtoy",
        	"@Version": "",
        	"TestData": [
        		{
        			"Id": [
        				1,
        				2,
        				3
        			]
        		}
        	]
        }`)
}

// 表头类型重复
func TestDuplicateKVHeaderType(t *testing.T) {

	emu := NewTableEmulator(t)

	dataSheet1 := emu.CreateDataSheet("data1", "")
	dataSheet1.WriteRow("Id")    // field name
	dataSheet1.WriteRow("int32") // field type
	dataSheet1.WriteRow()        // meta
	dataSheet1.WriteRow()        // comment

	// Key	Type	Value	Comment	Meta
	dataSheet2 := emu.CreateKVSheet("data2", "data1")
	dataSheet2.WriteRow("Id", "int32", "1", "", "")

	emu.MustGotError("TableError.DuplicateHeaderType 表头类型重复 | data1")
}
