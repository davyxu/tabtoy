package data

import (
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
)

// 动态protobuf message对象
type DynamicMessage struct {
	Desc *pbmeta.Descriptor

	fieldMap map[*pbmeta.FieldDescriptor]*fieldValue
}

// 字段值
type fieldValue struct {
	value      string
	valueArray []string

	msg      *DynamicMessage
	msgArray []*DynamicMessage
}

// 根据字段描述获取字段值,如果不存在自动创建
func (self *DynamicMessage) fetchValue(fd *pbmeta.FieldDescriptor, createIfNotExist bool) *fieldValue {
	fv, ok := self.fieldMap[fd]
	if !ok && createIfNotExist {
		fv = self.addValue(fd)
	}

	return fv
}

func (self *DynamicMessage) addValue(fd *pbmeta.FieldDescriptor) *fieldValue {
	fv := new(fieldValue)
	self.fieldMap[fd] = fv
	return fv
}

// 遍历所有值的描述符
func (self *DynamicMessage) IterateFieldDesc(callback func(*pbmeta.FieldDescriptor) bool) bool {

	for fd, _ := range self.fieldMap {
		if !callback(fd) {
			return false
		}
	}

	return true
}

// 值的描述符是否存在
func (self *DynamicMessage) ContainFieldDesc(fd *pbmeta.FieldDescriptor) bool {

	_, ok := self.fieldMap[fd]

	return ok
}

// 设置单一值
func (self *DynamicMessage) SetValue(fd *pbmeta.FieldDescriptor, value string) bool {

	if fd == nil || !self.Desc.Contains(fd) {
		log.Errorf("field not found: '%s' in '%s', value: '%s'", fd.Name(), self.Desc.Name(), value)
		return false
	}

	if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		log.Errorf("field is message not value: '%s' in '%s', value: '%s'", fd.Name(), self.Desc.Name(), value)
		return false
	}

	fv := self.fetchValue(fd, false)

	if fd.IsRepeated() {

		return false

	}

	// 可选字段
	if fd.IsOptional() {

		var existValue string
		if fv != nil {
			existValue = fv.value
		}

		// 有指派默认值
		if fd.DefaultValue() != "" {
			existValue = fd.DefaultValue()
		} else {
			// 没有指派默认值, 取值的默认值
			existValue = GetDefaultValue(fd)
		}

		// 输入值和已经存在的值一致, 就无需设置了
		if existValue == value {
			return true
		}

	}

	if fv == nil {
		fv = self.addValue(fd)
	}

	fv.value = value

	return true
}

// 添加值数组
func (self *DynamicMessage) AddRepeatedValue(fd *pbmeta.FieldDescriptor, value string) bool {

	if fd == nil || !self.Desc.Contains(fd) {
		log.Errorf("field not found: '%s' in '%s', value: '%s'", fd.Name(), self.Desc.Name(), value)
		return false
	}

	if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		log.Errorf("field is message not value: '%s' in '%s', value: '%s'", fd.Name(), self.Desc.Name(), value)
		return false
	}

	fv := self.fetchValue(fd, true)

	if !fd.IsRepeated() {
		return false
	}

	if fv.valueArray == nil {
		fv.valueArray = make([]string, 0)
	}

	fv.valueArray = append(fv.valueArray, value)

	return true
}

// 单一消息
func (self *DynamicMessage) SetMessage(fd *pbmeta.FieldDescriptor, value *DynamicMessage) bool {

	if fd == nil || !self.Desc.Contains(fd) {
		log.Errorf("field not found: '%s' in '%s'", fd.Name(), self.Desc.Name())
		return false
	}

	if fd.Type() != pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		log.Errorf("field is not message: '%s' in '%s'", fd.Name(), self.Desc.Name())
		return false
	}

	if fd.IsRepeated() {
		return false
	}

	self.fetchValue(fd, true).msg = value

	return true
}

// 添加消息数组
func (self *DynamicMessage) AddRepeatedMessage(fd *pbmeta.FieldDescriptor, value *DynamicMessage) bool {

	if fd == nil || !self.Desc.Contains(fd) {
		log.Errorf("field not found: '%s' in '%s'", fd.Name(), self.Desc.Name())
		return false
	}

	if fd.Type() != pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		log.Errorf("field is not message: '%s' in '%s'", fd.Name(), self.Desc.Name())
		return false
	}

	fv := self.fetchValue(fd, true)

	if !fd.IsRepeated() {
		return false
	}

	if fv.msgArray == nil {
		fv.msgArray = make([]*DynamicMessage, 0)
	}

	fv.msgArray = append(fv.msgArray, value)

	return true
}

// 取单一值
func (self *DynamicMessage) GetValue(fd *pbmeta.FieldDescriptor) (string, bool) {

	if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		return "", false
	}

	fv, ok := self.fieldMap[fd]

	if !ok {
		return "", false
	}

	if fd.IsRepeated() {
		return "", false
	}

	return fv.value, true
}

// 取值数组
func (self *DynamicMessage) GetRepeatedValue(fd *pbmeta.FieldDescriptor) []string {

	if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		return nil
	}

	fv, ok := self.fieldMap[fd]

	if !ok {
		return nil
	}

	if !fd.IsRepeated() {
		return nil
	}

	return fv.valueArray
}

// 取单一消息
func (self *DynamicMessage) GetMessage(fd *pbmeta.FieldDescriptor) *DynamicMessage {

	if fd.Type() != pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		return nil
	}

	fv, ok := self.fieldMap[fd]

	if !ok {
		return nil
	}

	if fd.IsRepeated() {
		return nil
	}

	return fv.msg
}

// 取消息数组
func (self *DynamicMessage) GetRepeatedMessage(fd *pbmeta.FieldDescriptor) []*DynamicMessage {

	if fd.Type() != pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		return nil
	}

	fv, ok := self.fieldMap[fd]

	if !ok {
		return nil
	}

	if !fd.IsRepeated() {
		return nil
	}

	return fv.msgArray
}

// 删除值
func (self *DynamicMessage) ClearFieldValue(fd *pbmeta.FieldDescriptor) bool {

	if fd == nil || !self.Desc.Contains(fd) {
		log.Errorf("field not found: '%s' in '%s' ", fd.Name(), self.Desc.Name())
		return false
	}

	if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {
		log.Errorf("field is message not value: '%s' in '%s'", fd.Name(), self.Desc.Name())
		return false
	}

	fv := self.fetchValue(fd, false)

	if fv == nil {
		return true
	}

	delete(self.fieldMap, fd)

	return true
}

func NewDynamicMessage(desc *pbmeta.Descriptor) *DynamicMessage {
	return &DynamicMessage{
		Desc:     desc,
		fieldMap: make(map[*pbmeta.FieldDescriptor]*fieldValue),
	}
}
