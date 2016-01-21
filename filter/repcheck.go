package filter

import (
	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/proto/tool"
)

type RepeatValueChecker map[string]bool

// 重复值检查， 如果meta中描述有开启则开启
func (self RepeatValueChecker) Check(meta *tool.FieldMeta, fd *pbmeta.FieldDescriptor, value string) {

	if meta == nil || meta.RepeatCheck == false {
		return
	}

	if self.contain(value) {
		log.Errorf("detected duplicate value %s=%s", fd.Name(), value)
		return
	}

	self[value] = true

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
