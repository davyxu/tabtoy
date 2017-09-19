package model

type Node struct {
	*FieldDescriptor

	StructRoot bool // 结构体标记的dummy node

	// 各种类型的值
	Value     string
	EnumValue int32
	Raw       []byte

	IValue interface{}

	Child []*Node // 优先遍历值, 再key

	SugguestIgnore bool //  建议忽略, 非repeated的普通字段导出时, 如果原单元格没填, 这个字段为true
}

func (self *Node) AddValue(value string, ivalue interface{}) *Node {

	n := &Node{
		Value:  value,
		IValue: ivalue,
	}
	self.Child = append(self.Child, n)

	return n
}

func (self *Node) AddKey(def *FieldDescriptor) *Node {
	n := &Node{
		FieldDescriptor: def,
	}
	self.Child = append(self.Child, n)

	return n
}
