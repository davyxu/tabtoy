package model

type FileDescriptor struct {
	Name             string
	DescriptorByName map[string]*Descriptor
	Descriptors      []*Descriptor

	Pragma *MetaInfo
}

func (self *FileDescriptor) MatchTag(tag string) bool {

	if !self.Pragma.ContainKey("OutputTag") {
		return true
	}

	return self.Pragma.ContainValue("OutputTag", tag)

}

// 取行类型的结构
func (self *FileDescriptor) RowDescriptor() *Descriptor {

	for _, d := range self.Descriptors {
		if d.Usage == DescriptorUsage_RowType {
			return d
		}
	}

	return nil
}

func (self *FileDescriptor) Add(def *Descriptor) bool {

	if _, ok := self.DescriptorByName[def.Name]; ok {

		return false
		//panic("duplicate buildin type")
	}

	// Descriptor会在每个表对应的FileDescriptor中和CombineFileDescriptor中同时存在
	// 这里忽略CombineFileDescriptor, 保持原有文件类型
	if def.File == nil {
		def.File = self
	}

	self.Descriptors = append(self.Descriptors, def)
	self.DescriptorByName[def.Name] = def

	return true
}

func NewFileDescriptor() *FileDescriptor {
	return &FileDescriptor{
		DescriptorByName: make(map[string]*Descriptor),
		Pragma:           NewMetaInfo(),
	}
}
