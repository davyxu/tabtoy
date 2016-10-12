package model

type DescriptorKind int

const (
	DescriptorKind_None DescriptorKind = iota
	DescriptorKind_Enum
	DescriptorKind_Struct
)

type DescriptorUsage int

const (
	DescriptorUsage_None          DescriptorUsage = iota
	DescriptorUsage_RowType                       // 每个表的行类型
	DescriptorUsage_CombineStruct                 // 最终使用的合并结构体
)

type Descriptor struct {
	Name  string
	Kind  DescriptorKind
	Usage DescriptorUsage

	FieldByName   map[string]*FieldDescriptor
	FieldByNumber map[int32]*FieldDescriptor
	Fields        []*FieldDescriptor

	Indexes     []*FieldDescriptor
	IndexByName map[string]*FieldDescriptor

	File *FileDescriptor
}

func (self *Descriptor) Add(def *FieldDescriptor) {

	def.Parent = self
	def.Order = int32(len(self.Fields))

	// 创建字段
	if _, ok := self.FieldByName[def.Name]; ok {
		panic("duplicate build in type")
		return
	} else {
		self.FieldByName[def.Name] = def
		self.FieldByNumber[def.EnumValue] = def
		self.Fields = append(self.Fields, def)
	}

	// 创建索引
	if def.Meta.MakeIndex {
		if _, ok := self.IndexByName[def.Name]; ok {
			panic("duplicate index name")
			return
		} else {
			self.IndexByName[def.Name] = def
			self.Indexes = append(self.Indexes, def)
		}
	}

}

func (self *Descriptor) FieldByValueAndMeta(value string) *FieldDescriptor {

	for _, v := range self.FieldByName {

		if v.Name == value {
			return v
		}

		if v.Meta.Alias == value {
			return v
		}

	}

	return nil
}

func NewDescriptor() *Descriptor {
	return &Descriptor{
		FieldByName:   make(map[string]*FieldDescriptor),
		FieldByNumber: make(map[int32]*FieldDescriptor),
		IndexByName:   make(map[string]*FieldDescriptor),
	}
}
