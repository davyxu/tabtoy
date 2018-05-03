package tests

import (
	"testing"
)

func TestTypeDefineOrder(t *testing.T) {

	mf := NewMemFile()
	indexSheet := mf.Create("Index.xlsx")

	WriteIndexTableHeader(indexSheet)
	WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")

	typeSheet := mf.Create("Type.xlsx")
	WriteTypeTableHeader(typeSheet)
	WriteRowValues(typeSheet, "表头", "TestData", "枚举", "Enum", "TestEnum", "", "")
	WriteRowValues(typeSheet, "枚举", "TestEnum", "", "None", "int", "", "0")

	if err := VerifyType(mf, ``); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestComplete(t *testing.T) {

	mf := NewMemFile()
	indexSheet := mf.Create("Index.xlsx")

	WriteIndexTableHeader(indexSheet)
	WriteRowValues(indexSheet, "类型表", "", "Type.xlsx")
	WriteRowValues(indexSheet, "数据表", "", "TestData.xlsx")

	typeSheet := mf.Create("Type.xlsx")
	WriteTypeTableHeader(typeSheet)
	WriteRowValues(typeSheet, "表头", "TestData", "整形", "Int", "int", "", "")
	WriteRowValues(typeSheet, "表头", "TestData", "字符串", "String", "string", "", "")
	WriteRowValues(typeSheet, "表头", "TestData", "布尔", "Bool", "bool", "", "")

	dataSheet := mf.Create("TestData.xlsx")
	WriteRowValues(dataSheet, "整形", "字符串", "布尔")
	WriteRowValues(dataSheet, "100", "\"hello\"", "true")

	if err := VerifyLauncherJson(mf, `
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
