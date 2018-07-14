package table

type TableKind int32

const (
	TableKind_None         = 0 //
	TableKind_HeaderStruct = 1 // 表头
	TableKind_Enum         = 2 // 枚举
)

func (self TableKind) String() string {

	switch self {
	case TableKind_HeaderStruct:
		return "表头"
	case TableKind_Enum:
		return "枚举"
	default:
		return "未知"
	}
}

type TableField struct {
	Kind          TableKind `tb_name:"种类"`
	ObjectType    string    `tb_name:"对象类型"`
	Name          string    `tb_name:"标识名"`
	FieldName     string    `tb_name:"字段名"`
	FieldType     string    `tb_name:"字段类型"`
	Value         string    `tb_name:"值" json:",omitempty"`
	ArraySplitter string    `tb_name:"数组切割" json:",omitempty"`
	MakeIndex     bool      `tb_name:"索引" json:",omitempty"`
	IsBuiltin     bool      `json:",omitempty"`
}

func (self *TableField) IsArray() bool {
	return self.ArraySplitter != ""
}

type TableMode int32

const (
	TableMode_None     = 0 //
	TableMode_Type     = 1 // 类型表
	TableMode_Data     = 2 // 数据表
	TableMode_KeyValue = 3 // 键值表
)

var (
	TableModeMapperValueByName = map[string]int32{
		"None":     0, //
		"Type":     1, // 类型表
		"Data":     2, // 数据表
		"KeyValue": 3, // 键值表
	}

	TableModeMapperNameByValue = map[int32]string{
		0: "None",     //
		1: "Type",     // 类型表
		2: "Data",     // 数据表
		3: "KeyValue", // 键值表
	}
)

type TablePragma struct {
	TableMode     TableMode `tb_name:"模式"`
	TableType     string    `tb_name:"表类型"`
	TableFileName string    `tb_name:"表文件名"`
}
