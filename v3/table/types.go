package table

type TableKind int32

const (
	TableKind_None         = 0 //
	TableKind_HeaderStruct = 1 // 表头
	TableKind_Enum         = 2 // 枚举
)

type TableField struct {
	Kind          TableKind `tb_name:"种类"`
	ObjectType    string    `tb_name:"对象类型"`
	Name          string    `tb_name:"标识名"`
	FieldName     string    `tb_name:"字段名"`
	FieldType     string    `tb_name:"字段类型"`
	Value         string    `tb_name:"值" json:",omitempty"`
	ArraySplitter string    `tb_name:"数组切割" json:",omitempty"`
	IsBuiltin     bool      `json:",omitempty"`
}

func (self *TableField) IsArray() bool {
	return self.ArraySplitter != ""
}
