package model

type Node struct {
	*FieldDefine

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

func (self *Node) AddKey(def *FieldDefine) *Node {
	n := &Node{
		FieldDefine: def,
	}
	self.Child = append(self.Child, n)

	return n
}

type Record struct {
	nodeByFD map[*FieldDefine]*Node
	Nodes    []*Node
}

func (self *Record) NewNodeByDefine(def *FieldDefine) *Node {

	// 如果这个单元格数据有, 使用已经有的定义, 因为字段不会重复
	// 主要处理repeated散开的case
	if exist, ok := self.nodeByFD[def]; ok {
		return exist
	}

	node := new(Node)
	node.FieldDefine = def
	self.nodeByFD[def] = node
	self.Nodes = append(self.Nodes, node)

	return node
}

func NewRecord() *Record {
	return &Record{
		nodeByFD: make(map[*FieldDefine]*Node),
	}
}
