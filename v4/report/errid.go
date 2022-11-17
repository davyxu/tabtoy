package report

import "github.com/davyxu/tabtoy/util"

func Init() {
	util.RegisterError("InvalidMetaKey", &util.ErrorLanguage{CHS: "非法的字段元属性键"})
	util.RegisterError("InvalidKVHeader", &util.ErrorLanguage{CHS: "非法KV表头"})
	util.RegisterError("DuplicateKVField", &util.ErrorLanguage{CHS: "键值表字段重复"})
	util.RegisterError("UnknownFieldType", &util.ErrorLanguage{CHS: "未知字段类型"})
	//util.RegisterError("DuplicateTypeFieldName", &util.ErrorLanguage{CHS: "类型表字段重复"})
	//util.RegisterError("EnumValueEmpty", &util.ErrorLanguage{CHS: "枚举值空"})
	//util.RegisterError("DuplicateEnumValue", &util.ErrorLanguage{CHS: "枚举值重复"})
	//util.RegisterError("UnknownEnumValue", &util.ErrorLanguage{CHS: "未知的枚举值"})
	//util.RegisterError("InvalidTypeTable", &util.ErrorLanguage{CHS: "非法的类型表"})
	//util.RegisterError("DuplicateValueInMakingIndex", &util.ErrorLanguage{CHS: "创建索引时发现重复值"})
	//util.RegisterError("UnknownInputFileExtension", &util.ErrorLanguage{CHS: "未知的输入文件扩展名"})
	//util.RegisterError("DataMissMatchTypeDefine", &util.ErrorLanguage{CHS: "数据与定义类型不匹配"})
	//util.RegisterError("ArrayMultiColumnDefineNotMatch", &util.ErrorLanguage{CHS: "数组类型多列跨表定义不一致"})
	//util.RegisterError("InvalidFieldName", &util.ErrorLanguage{CHS: "非法字段名"})
	//util.RegisterError("UnknownTypeKind", &util.ErrorLanguage{CHS: "非法的类型种类"})
}
