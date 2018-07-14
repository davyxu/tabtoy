package table

import (
	"github.com/davyxu/tabtoy/v3/helper"
)

var BuiltinTypes helper.FileGetter

func createBuiltinTypes(mf *helper.MemFile) {

	sheet := mf.CreateDefault("BuiltinTypes.xlsx")
	helper.WriteTypeTableHeader(sheet)
	// "种类", "对象类型", "标识名", "字段名", "字段类型", "数组切割", "值", "索引"

	// 索引表格式
	helper.WriteRowValues(sheet, "表头", "TablePragma", "模式", "TableMode", "TableMode", "", "")
	helper.WriteRowValues(sheet, "表头", "TablePragma", "表类型", "TableType", "string", "", "")
	helper.WriteRowValues(sheet, "表头", "TablePragma", "表文件名", "TableFileName", "string", "", "")

	// 索引表类型
	helper.WriteRowValues(sheet, "枚举", "TableMode", "", "None", "int", "", "0")
	helper.WriteRowValues(sheet, "枚举", "TableMode", "类型表", "Type", "int", "", "1")
	helper.WriteRowValues(sheet, "枚举", "TableMode", "数据表", "Data", "int", "", "2")
	helper.WriteRowValues(sheet, "枚举", "TableMode", "键值表", "KeyValue", "int", "", "3")

	// KV表格式
	helper.WriteRowValues(sheet, "表头", "TableKeyValue", "字段名", "FieldName", "string", "", "")
	helper.WriteRowValues(sheet, "表头", "TableKeyValue", "字段类型", "FieldType", "string", "", "")
	helper.WriteRowValues(sheet, "表头", "TableKeyValue", "标识名", "Name", "string", "", "")
	helper.WriteRowValues(sheet, "表头", "TableKeyValue", "值", "Value", "string", "", "")
	helper.WriteRowValues(sheet, "表头", "TableKeyValue", "数组切割", "ArraySplitter", "string", "", "")
}

func init() {
	memFile := helper.NewMemFile()

	createBuiltinTypes(memFile)

	BuiltinTypes = memFile
}
