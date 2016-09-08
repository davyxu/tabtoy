package filter

import (
	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/proto/tool"
)

type RepeatValueChecker map[string]bool

// 重复值检查， 如果meta中描述有开启则开启
func (self RepeatValueChecker) Check(meta *tool.FieldMetaV1, fd *pbmeta.FieldDescriptor, value string) bool {

	if meta == nil || meta.RepeatCheck == false {
		return true
	}

	if self.contain(value) {
		log.Errorf("detected duplicate value %s=%s", fd.Name(), value)
		return false
	}

	self[value] = true

	return true
}

func (self RepeatValueChecker) contain(value string) bool {

	if _, ok := self[value]; ok {
		return true
	}

	return false
}

func NewRepeatValueChecker() RepeatValueChecker {
	return make(map[string]bool)
}
