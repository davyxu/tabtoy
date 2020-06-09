package jsontype

import "strconv"

type Field struct {
	Name          string // 字段名
	Type          string // 表中原有写的类型
	Comment       string // 注释, 表中的名称
	Value         string `json:",omitempty"` // 枚举值
	MakeIndex     bool   `json:",omitempty"` // 是否生成索引
	ArraySplitter string `json:",omitempty"` // 数组切割符
}

func (self *Field) EnumValue() int {
	v, _ := strconv.Atoi(self.Value)
	return v
}

// 表示表格, 枚举
type Object struct {
	Name string
	Type string // 对象类型

	Fields []*Field
}

func (self *Object) Compare(other *Object) bool {
	if self.Type != other.Type {
		return self.Type < other.Type
	}

	return self.Name < other.Name
}

type File struct {
	Tool    string `json:"@Tool"`
	Version string `json:"@Version"`
	Objects []*Object
}
