package report

import "github.com/davyxu/tabtoy/util"

func Init() {
	util.RegisterError("UnknownMetaKey", &util.ErrorLanguage{CHS: "未知的元属性键"})
	util.RegisterError("EmptyArraySpliter", &util.ErrorLanguage{CHS: "数组切割符空"})
	util.RegisterError("InvalidKVHeader", &util.ErrorLanguage{CHS: "非法KV表头"})
	util.RegisterError("DuplicateKVField", &util.ErrorLanguage{CHS: "键值表字段重复"})
	util.RegisterError("UnknownFieldType", &util.ErrorLanguage{CHS: "未知字段类型"})
	util.RegisterError("UnknownFieldName", &util.ErrorLanguage{CHS: "未知字段名"})
	util.RegisterError("DuplicateHeaderField", &util.ErrorLanguage{CHS: "表头字段重复"})
	util.RegisterError("InvalidMetaFormat", &util.ErrorLanguage{CHS: "非法的元格式"})
	util.RegisterError("DuplicateValueInMakingIndex", &util.ErrorLanguage{CHS: "创建索引时发现重复值"})
}
