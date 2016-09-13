package model

type Record struct {
	nodeByFD map[*FieldDescriptor]*Node
	Nodes    []*Node
}

func (self *Record) NewNodeByDefine(def *FieldDescriptor) *Node {

	// 如果这个单元格数据有, 使用已经有的定义, 因为字段不会重复
	// 主要处理repeated散开的case
	if exist, ok := self.nodeByFD[def]; ok {
		return exist
	}

	node := new(Node)
	node.FieldDescriptor = def
	self.nodeByFD[def] = node
	self.Nodes = append(self.Nodes, node)

	return node
}

func NewRecord() *Record {
	return &Record{
		nodeByFD: make(map[*FieldDescriptor]*Node),
	}
}

type Table struct {
	*FileDescriptor
	Recs []*Record
}

func (self *Table) Add(r *Record) {
	self.Recs = append(self.Recs, r)
}

func NewTable(fileD *FileDescriptor) *Table {
	return &Table{
		FileDescriptor: fileD,
	}
}
