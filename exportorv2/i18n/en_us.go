package i18n

func init() {

	registerLanguage("en_us", map[StringID]string{
		ConvertValue_EnumTypeNil:                         "ConvertValue: Enum type nil",
		ConvertValue_StructTypeNil:                       "ConvertValue: Struct type nil",
		ConvertValue_EnumValueNotFound:                   "ConvertValue: Enum value not found",
		ConvertValue_UnknownFieldType:                    "ConvertValue: Unknown field type",
		StructParser_LexerError:                          "StructParser: Lexer error",
		StructParser_ExpectField:                         "StructParser: Expect field",
		StructParser_UnexpectedSpliter:                   "StructParser: Unexpected k-v spliter",
		StructParser_FieldNotFound:                       "StructParser: Field not found",
		StructParser_DuplicateFieldInCell:                "StructParser: Duplicate field",
		Run_CollectTypeInfo:                              "Run: Collect Type Info",
		Run_ExportSheetData:                              "Run: Export Sheet Data",
		Globals_CombineNameLost:                          "Globals: Please specify 'combinename' params",
		Globals_PackageNameDiff:                          "Globals: Keep all type in same package",
		Globals_TableNameDuplicated:                      "Globals: Duplicate table name",
		Globals_OutputCombineData:                        "Globals: Merge Combined Data",
		File_TypeSheetKeepSingleton:                      "File: Type sheet only need ONE in a file",
		DataSheet_ValueConvertError:                      "DataSheet: Cell value convert error",
		DataSheet_ValueRepeated:                          "DataSheet: Duplicated cell value",
		DataHeader_TypeNotFound:                          "DataHeader: Type not found",
		DataHeader_MetaParseFailed:                       "DataHeader: Meta parse failed",
		DataHeader_DuplicateFieldName:                    "DataHeader: Duplicated field name",
		DataHeader_RepeatedFieldTypeNotSameInMultiColumn: "DataHeader: Repeated field type not same in columns",
		DataHeader_RepeatedFieldMetaNotSameInMultiColumn: "DataHeader: Repeated field meta not same in columns",
		TypeSheet_PragmaParseFailed:                      "TypeSheet: File pragma parse failed",
		TypeSheet_TableNameIsEmpty:                       "TypeSheet: Table name is empty",
		TypeSheet_PackageIsEmpty:                         "TypeSheet: Package is empty",
		TypeSheet_FieldTypeNotFound:                      "TypeSheet: Field type not found",
		TypeSheet_EnumValueParseFailed:                   "TypeSheet: Enum value parse failed",
		TypeSheet_DescriptorKindNotSame:                  "TypeSheet: Descriptor kind not the same, due to enum value",
		TypeSheet_FieldMetaParseFailed:                   "TypeSheet: Field meta parse failed",
		TypeSheet_StructFieldCanNotBeStruct:              "TypeSheet: Struct field can not be struct kind",
		TypeSheet_FirstEnumValueShouldBeZero:             "TypeSheet: First enum value should be zero",
	})
}
