package model

type Node struct {
	*FieldDescriptor

	StructRoot bool // 结构体标记的dummy node

	// 各种类型的值
	Value     string
	EnumValue int32
	Raw       []byte

	Child []*Node // 优先遍历值, 再key

	SugguestIgnore bool //  建议忽略, 非repeated的普通字段导出时, 如果原单元格没填, 这个字段为true
}

func (self *Node) AddValue(value string) *Node {

	n := &Node{
		Value: value,
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
