package model

import (
	"github.com/davyxu/tabtoy/proto/tool"
)

type FileDescriptor struct {
	Name             string
	DescriptorByName map[string]*Descriptor
	Descriptors      []*Descriptor

	FileType *Descriptor // 自动创建的XXFile类型, 一个BuildInTypeSet 一次只有1个这样的对象

	Pragma tool.FilePragmaV2
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

func (self *FileDescriptor) Add(def *Descriptor) {

	if _, ok := self.DescriptorByName[def.Name]; ok {
		panic("duplicate buildin type")
	}

	def.File = self
	self.Descriptors = append(self.Descriptors, def)
	self.DescriptorByName[def.Name] = def
}

func NewFileDescriptor() *FileDescriptor {
	return &FileDescriptor{
		DescriptorByName: make(map[string]*Descriptor),
	}
}
