package printer

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
)

type jsonWriter struct {
	printer *bytes.Buffer
}

func (self *jsonWriter) RepeatedMessageBegin(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int) {

	if indent == 1 {
		self.printer.WriteString(fmt.Sprintf("\"%s\" : [\n", fd.Name()))
	} else {
		self.printer.WriteString(fmt.Sprintf("\"%s\" : [", fd.Name()))
	}

}

// Value是消息的字段
func (self *jsonWriter) WriteMessage(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, indent int) {

	pos := self.printer.Len()

	if indent == 1 || fd.IsRepeated() {
		self.printer.WriteString("{")

	} else {
		self.printer.WriteString(fmt.Sprintf("\"%s\" : {", fd.Name()))
	}

	valueWrite := rawWriteMessage(self.printer, self, msg, indent)

	self.printer.WriteString("}")

	// 如果没写入值, 将回滚写入
	if valueWrite == 0 {
		self.printer.Truncate(pos)
	}
}

func (self *jsonWriter) RepeatedMessageEnd(fd *pbmeta.FieldDescriptor, msg *data.DynamicMessage, msgCount int, indent int) {

	self.printer.WriteString("]")

}

func (self *jsonWriter) RepeatedValueBegin(fd *pbmeta.FieldDescriptor) {
	self.printer.WriteString(fmt.Sprintf("\"%s\" : [", fd.Name()))
}

func jsonvalueWrapper(fd *pbmeta.FieldDescriptor, value string) string {

	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_STRING:
		return strEscape(value)
	case pbprotos.FieldDescriptorProto_TYPE_ENUM:
		ed := fd.EnumDesc()

		if ed == nil {
			return "0"
		}

		evd := ed.ValueByName(value)

		if evd == nil {
			return "0"
		}

		return strconv.Itoa(int(evd.Value()))

	}

	return value
}

// 普通值
func (self *jsonWriter) WriteValue(fd *pbmeta.FieldDescriptor, value string, indent int) {

	if fd.IsRepeated() {
		self.printer.WriteString(jsonvalueWrapper(fd, value))
	} else {
		self.printer.WriteString(fmt.Sprintf("\"%s\": %s", fd.Name(), jsonvalueWrapper(fd, value)))
	}

}

func (self *jsonWriter) RepeatedValueEnd(fd *pbmeta.FieldDescriptor) {
	self.printer.WriteString("]")
}

func (self *jsonWriter) WriteFieldSpliter() {

	self.printer.WriteString(", ")
}

func (self *jsonWriter) PrintMessage(msg *data.DynamicMessage) bool {

	self.printer.WriteString("{\n\n")
	rawWriteMessage(self.printer, self, msg, 0)

	self.printer.WriteString("\n\n}")

	return true
}

func NewJsonWriter(printer *bytes.Buffer) IWriter {

	return &jsonWriter{
		printer: printer,
	}
}
