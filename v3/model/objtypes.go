package model

// 类型表中的一行描述
type TypeField struct {
	Table        string `tab:"表名"`
	ObjectType   string `tab:"对象类型"`
	Name         string `tab:"标识名"`
	FieldName    string `tab:"字段名"`
	FieldType    string `tab:"字段类型"`
	DefaultValue string `tab:"默认值"` // 枚举值
}
