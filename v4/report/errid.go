package report

import "github.com/davyxu/tabtoy/util"

func Init() {
	util.RegisterError("UnknownMetaKey", &util.ErrorLanguage{CHS: "未知的元属性键"})
	util.RegisterError("EmptyArraySpliter", &util.ErrorLanguage{CHS: "数组切割符空"})
	util.RegisterError("InvalidKVHeader", &util.ErrorLanguage{CHS: "非法键值表头"})
	util.RegisterError("DuplicateKVField", &util.ErrorLanguage{CHS: "KV表字段重复"})
	util.RegisterError("UnknownFieldType", &util.ErrorLanguage{CHS: "未知字段类型"})
	util.RegisterError("UnknownFieldName", &util.ErrorLanguage{CHS: "未知字段名"})
	util.RegisterError("DuplicateHeaderField", &util.ErrorLanguage{CHS: "表头字段重复"})
	util.RegisterError("InvalidMetaFormat", &util.ErrorLanguage{CHS: "非法的元格式"})
	util.RegisterError("DuplicateValueInMakingIndex", &util.ErrorLanguage{CHS: "创建索引时发现重复值"})
	util.RegisterError("DuplicateHeaderType", &util.ErrorLanguage{CHS: "表头类型重复"})
	util.RegisterError("InvalidIndexHeader", &util.ErrorLanguage{CHS: "非法索引表表头"})
	util.RegisterError("EmptyTableType", &util.ErrorLanguage{CHS: "索引表中未指定表类型"})
	util.RegisterError("InvalidTableMode", &util.ErrorLanguage{CHS: "索引表中未指定表模式"})
	util.RegisterError("InvalidTypeHeader", &util.ErrorLanguage{CHS: "非法类型表头"})
	util.RegisterError("DuplicateTypeField", &util.ErrorLanguage{CHS: "类型表字段重复"})
	util.RegisterError("RequireInteger", &util.ErrorLanguage{CHS: "类型表需要整数值"})
	util.RegisterError("DuplicateEnumValue", &util.ErrorLanguage{CHS: "枚举值重复"})
	util.RegisterError("EnumValueEmpty", &util.ErrorLanguage{CHS: "枚举值空"})
	util.RegisterError("UnknownEnumValue", &util.ErrorLanguage{CHS: "未知的枚举值"})
}
