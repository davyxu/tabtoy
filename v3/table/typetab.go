package table

import (
	"github.com/davyxu/tabtoy/v3/helper"
)

var BuiltinTypes helper.FileGetter

func init() {
	mf := helper.NewMemFile()

	// TODO 支持导出这种代码
	typeSheet := mf.Create("BuiltinTypes.xlsx")
	helper.WriteTypeTableHeader(typeSheet)
	helper.WriteRowValues(typeSheet, "表头", "FieldType", "输入字段", "InputFieldName", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "FieldType", "Go字段", "GoFieldName", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "FieldType", "C#字段", "CSFieldName", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "FieldType", "默认值", "DefaultValue", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TablePragma", "模式", "TableMode", "TableMode", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TablePragma", "表类型", "TableType", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TablePragma", "表文件名", "TableFileName", "string", "", "")
	helper.WriteRowValues(typeSheet, "枚举", "TableMode", "", "None", "int", "", "0")
	helper.WriteRowValues(typeSheet, "枚举", "TableMode", "数据表", "Data", "int", "", "1")
	helper.WriteRowValues(typeSheet, "枚举", "TableMode", "类型表", "Type", "int", "", "2")
	helper.WriteRowValues(typeSheet, "枚举", "TableMode", "键值表", "KeyValue", "int", "", "3")
	helper.WriteRowValues(typeSheet, "表头", "TableKeyValue", "字段名", "FieldName", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TableKeyValue", "字段类型", "FieldType", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TableKeyValue", "标识名", "Name", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TableKeyValue", "值", "Value", "string", "", "")
	helper.WriteRowValues(typeSheet, "表头", "TableKeyValue", "数组切割", "ArraySplitter", "string", "", "")

	BuiltinTypes = mf
}
