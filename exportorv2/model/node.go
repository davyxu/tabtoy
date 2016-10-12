package model

type Node struct {
	*FieldDescriptor

	StructRoot bool // 结构体标记的dummy node

	// 各种类型的值
	Value     string
	EnumValue int32
	Raw       []byte

	Child []*Node // 优先遍历值, 再key
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
