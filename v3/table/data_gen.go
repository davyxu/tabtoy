package table

const coreConfig = `
{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "3.0.0",	
	"FieldType":[ 
		{ "InputFieldName": "int32", "GoFieldName": "int32", "CSFieldName": "Int32", "DefaultValue": "0" },
		{ "InputFieldName": "int64", "GoFieldName": "int64", "CSFieldName": "Int64", "DefaultValue": "0" },
		{ "InputFieldName": "int", "GoFieldName": "int32", "CSFieldName": "Int32", "DefaultValue": "0" },
		{ "InputFieldName": "uint32", "GoFieldName": "uint32", "CSFieldName": "UInt32", "DefaultValue": "0" },
		{ "InputFieldName": "uint64", "GoFieldName": "uint64", "CSFieldName": "UInt64", "DefaultValue": "0" },
		{ "InputFieldName": "float", "GoFieldName": "float32", "CSFieldName": "float", "DefaultValue": "0" },
		{ "InputFieldName": "double", "GoFieldName": "float64", "CSFieldName": "double", "DefaultValue": "0" },
		{ "InputFieldName": "float32", "GoFieldName": "float32", "CSFieldName": "float", "DefaultValue": "0" },
		{ "InputFieldName": "float64", "GoFieldName": "float64", "CSFieldName": "double", "DefaultValue": "0" },
		{ "InputFieldName": "bool", "GoFieldName": "bool", "CSFieldName": "bool", "DefaultValue": "false" },
		{ "InputFieldName": "string", "GoFieldName": "string", "CSFieldName": "string", "DefaultValue": "" } 
	],
	"ErrorID":[ 
		{ "HeaderNotMatchFieldName": "表头与字段不匹配", "HeaderFieldNotDefined": "表头字段未定义", "DuplicateHeaderField": "表头字段重复", "DuplicateKVField": "键值表字段重复", "UnknownFieldType": "未知字段类型", "DuplicateTypeFieldName": "类型表字段重复", "EnumValueEmpty": "枚举值空", "DuplicateEnumValue": "枚举值重复" } 
	]
}`
