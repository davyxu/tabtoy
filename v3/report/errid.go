package report

type ErrorLanguage struct {
	CHS string
}

var (
	ErrorByID = map[string]*ErrorLanguage{
		"HeaderNotMatchFieldName":     {CHS: "表头与字段不匹配"},
		"HeaderFieldNotDefined":       {CHS: "表头字段未定义"},
		"DuplicateHeaderField":        {CHS: "表头字段重复"},
		"DuplicateKVField":            {CHS: "键值表字段重复"},
		"UnknownFieldType":            {CHS: "未知字段类型"},
		"DuplicateTypeFieldName":      {CHS: "类型表字段重复"},
		"EnumValueEmpty":              {CHS: "枚举值空"},
		"DuplicateEnumValue":          {CHS: "枚举值重复"},
		"UnknownEnumValue":            {CHS: "未知的枚举值"},
		"InvalidTypeTable":            {CHS: "非法的类型表"},
		"HeaderTypeNotFound":          {CHS: "表头类型找不到"},
		"DuplicateValueInMakingIndex": {CHS: "创建索引时发现重复值"},
		"UnknownInputFileExtension":   {CHS: "未知的输入文件扩展名"},
		"DataMissMatchTypeDefine":     {CHS: "数据与定义类型不匹配"},
	}
)
