package tests

import "testing"

// 类型字段名重复
func TestTypeDuplicateFieldName(t *testing.T) {

	emu := NewTableEmulator(t)

	typeSheet := emu.CreateTypeSheet(TestCSVName)
	// ObjectType	FieldName	Value	Comment
	typeSheet.WriteRow("ItemType", "A", "1", "")
	typeSheet.WriteRow("ItemType", "A", "1", "")

	emu.MustGotError("TableError.DuplicateTypeField 类型表字段重复 | 'A' @TestData|(B3)")
}

// 类型值重复
func TestTypeDuplicateValue(t *testing.T) {

	emu := NewTableEmulator(t)

	typeSheet := emu.CreateTypeSheet(TestCSVName)
	// ObjectType	FieldName	Value	Comment
	typeSheet.WriteRow("ItemType", "A", "1", "")
	typeSheet.WriteRow("ItemType", "B", "1", "")

	emu.MustGotError("TableError.DuplicateEnumValue 枚举值重复 | 1 'B' @TestData|(C3)")
}

// 类型值空
func TestTypeEmptyValue(t *testing.T) {

	emu := NewTableEmulator(t)

	typeSheet := emu.CreateTypeSheet(TestCSVName)
	// ObjectType	FieldName	Value	Comment
	typeSheet.WriteRow("ItemType", "A", "", "")

	emu.MustGotError("TableError.EnumValueEmpty 枚举值空 |  'A' @TestData|(C2)")
}

// 枚举值使用时找不到
func TestTypeUnknownEnumValue(t *testing.T) {

	emu := NewTableEmulator(t)
	dataSheet := emu.CreateDataSheet(TestCSVName, "")
	dataSheet.WriteRow("Type")     // field name
	dataSheet.WriteRow("ItemType") // field type
	dataSheet.WriteRow("")         // meta
	dataSheet.WriteRow()           // comment
	dataSheet.WriteRow("B")

	typeSheet := emu.CreateTypeSheet("Type")
	// ObjectType	FieldName	Value	Comment
	typeSheet.WriteRow("ItemType", "A", "1", "")

	emu.MustGotError("TableError.UnknownEnumValue 未知的枚举值 | ItemType 'B' @TestData|TestData(A1)")
}
