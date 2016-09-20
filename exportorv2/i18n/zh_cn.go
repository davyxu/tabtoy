package i18n

func init() {

	registerLanguage("zh_cn", map[StringID]string{
		ConvertValue_EnumTypeNil:                         "值转换: 枚举类型空",
		ConvertValue_StructTypeNil:                       "值转换: 结构类型空",
		ConvertValue_EnumValueNotFound:                   "值转换: 枚举值未找到",
		ConvertValue_UnknownFieldType:                    "值转换: 未知的字段类型",
		StructParser_LexerError:                          "结构体解析: 词法错误",
		StructParser_ExpectField:                         "结构体解析: 期望字段",
		StructParser_UnexpectedSpliter:                   "结构体解析: 非预期的键值分割符",
		StructParser_FieldNotFound:                       "结构体解析: 未知字段",
		Run_CollectTypeInfo:                              "运行: 收集类型信息",
		Run_ExportSheetData:                              "运行: 导出表单数据",
		Globals_CombineNameLost:                          "合并: 请在参数中添加 'combinename' 指明合并配置名",
		Globals_PackageNameDiff:                          "合并: 所有表中的@Types中的包名(Package)请保持一致",
		Globals_TableNameDuplicated:                      "合并: 表名(TableName)重复",
		Globals_OutputCombineData:                        "合并: 输出合并数据",
		File_TypeSheetKeepSingleton:                      "文件: 类型表在一个表中只能有一份",
		DataSheet_ValueConvertError:                      "数据表: 单元格转换错误",
		DataSheet_ValueRepeated:                          "数据表: 单元格值重复",
		DataHeader_TypeNotFound:                          "数据头: 未知类型",
		DataHeader_MetaParseFailed:                       "数据头: 特性解析错误",
		DataHeader_DuplicateFieldName:                    "数据头: 重复字段名",
		DataHeader_RepeatedFieldTypeNotSameInMultiColumn: "数据头: 数组字段在多列中的类型不一致",
		DataHeader_RepeatedFieldMetaNotSameInMultiColumn: "数据头: 数组字段在多列中的特性不一致",
		TypeSheet_PragmaParseFailed:                      "类型表: 文件特性解析失败",
		TypeSheet_TableNameIsEmpty:                       "类型表: 表名(TableName)为空",
		TypeSheet_PackageIsEmpty:                         "类型表: 包名(Package)为空",
		TypeSheet_FieldTypeNotFound:                      "类型表: 未知字段类型",
		TypeSheet_EnumValueParseFailed:                   "类型表: 枚举值解析失败",
		TypeSheet_DescriptorKindNotSame:                  "类型表: 类型前后不一致, 由枚举值不一致导致",
		TypeSheet_FieldMetaParseFailed:                   "类型表: 字段特性解析失败",
		TypeSheet_StructFieldCanNotBeStruct:              "类型表: 结构体字段类型不能是结构体类型",
		TypeSheet_FirstEnumValueShouldBeZero:             "类型表: 第一个枚举值必须为0",
	})
}
