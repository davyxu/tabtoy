package model

func InitBuiltinTypes(typeTab *TypeTable) {

	for _, tf := range []*TypeDefine{

		// 类型表类型
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "", FieldName: "None", FieldType: "int", Value: "0"},
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "表头", FieldName: "HeaderStruct", FieldType: "int", Value: "1"},
		{Kind: TypeUsage_Enum, ObjectType: "TypeUsage", Name: "枚举", FieldName: "Enum", FieldType: "int", Value: "2"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "种类", FieldName: "Kind", FieldType: "TypeUsage"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "对象类型", FieldName: "ObjectType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "标识名", FieldName: "Name", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "字段名", FieldName: "FieldName", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "字段类型", FieldName: "FieldType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "值", FieldName: "Value", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "TypeDefine", Name: "数组切割", FieldName: "ArraySplitter", FieldType: "string"},

		// 索引表类型
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "", FieldName: "None", FieldType: "int", Value: "0"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "类型表", FieldName: "Type", FieldType: "int", Value: "1"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "数据表", FieldName: "Data", FieldType: "int", Value: "2"},
		{Kind: TypeUsage_Enum, ObjectType: "TableKind", Name: "键值表", FieldName: "KeyValue", FieldType: "int", Value: "3"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "模式", FieldName: "TableKind", FieldType: "TableKind"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "表类型", FieldName: "TableType", FieldType: "TableKind"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "IndexDefine", Name: "表文件名", FieldName: "TableFileName", FieldType: "TableKind"},

		// KV表类型
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "字段名", FieldName: "FieldName", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "字段类型", FieldName: "FieldType", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "标识名", FieldName: "Name", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "值", FieldName: "Value", FieldType: "string"},
		{Kind: TypeUsage_HeaderStruct, ObjectType: "KVDefine", Name: "数组切割", FieldName: "ArraySplitter", FieldType: "string"},
	} {
		tf.IsBuiltin = true

		typeTab.AddField(tf, nil, 0)
	}
}
