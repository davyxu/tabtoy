package model

type Node struct {
	Define *FieldDefine
	Value  string

	StructRoot bool

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
		Define: def,
	}
	self.Child = append(self.Child, n)

	return n
}

type Record struct {
	nodeByFD map[*FieldDefine]*Node
	nodes    []*Node
}

func (self *Record) NewNodeByDefine(def *FieldDefine) *Node {

	// 如果这个单元格数据有, 使用已经有的定义, 因为字段不会重复
	// 主要处理repeated散开的case
	if exist, ok := self.nodeByFD[def]; ok {
		return exist
	}

	node := new(Node)
	node.Define = def
	self.nodeByFD[def] = node
	self.nodes = append(self.nodes, node)

	return node
}

func NewRecord() *Record {
	return &Record{
		nodeByFD: make(map[*FieldDefine]*Node),
	}
}
