package printer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
	"os"
	"path/filepath"
)

type Stream struct {
	buf bytes.Buffer
}

func (self *Stream) Len() int {
	return self.buf.Len()
}

func (self *Stream) Buffer() *bytes.Buffer {
	return &self.buf
}

func (self *Stream) WriteBytes(b []byte) {
	self.buf.Write(b)
}

func (self *Stream) Printf(format string, args ...interface{}) {
	self.buf.WriteString(fmt.Sprintf(format, args...))
}

func (self *Stream) WriteFile(outfile string) error {

	// 自动创建目录
	os.MkdirAll(filepath.Dir(outfile), 666)

	err := ioutil.WriteFile(outfile, self.buf.Bytes(), 0666)
	if err != nil {
		log.Errorf("%s, %v", i18n.String(i18n.Printer_OpenWriteOutputFileFailed), err.Error())
		return err
	}

	return nil
}

func (self *Stream) WriteInt32(v int32) {

	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *Stream) WriteString(v string) {
	rawStr := []byte(v)

	binary.Write(&self.buf, binary.LittleEndian, int32(len(rawStr)))

	binary.Write(&self.buf, binary.LittleEndian, rawStr)
}

func (self *Stream) WriteNodeValue(ft model.FieldType, value *model.Node) {

	switch ft {
	case model.FieldType_Int32:
		v, _ := strconv.ParseInt(value.Value, 10, 32)
		binary.Write(&self.buf, binary.LittleEndian, int32(v))
	case model.FieldType_UInt32:
		v, _ := strconv.ParseUint(value.Value, 10, 32)

		binary.Write(&self.buf, binary.LittleEndian, uint32(v))
	case model.FieldType_Int64:
		v, _ := strconv.ParseInt(value.Value, 10, 64)

		binary.Write(&self.buf, binary.LittleEndian, int64(v))
	case model.FieldType_UInt64:
		v, _ := strconv.ParseUint(value.Value, 10, 64)

		binary.Write(&self.buf, binary.LittleEndian, uint64(v))
	case model.FieldType_Float:
		v, _ := strconv.ParseFloat(value.Value, 32)

		binary.Write(&self.buf, binary.LittleEndian, float32(v))
	case model.FieldType_Bool:
		v, _ := strconv.ParseBool(value.Value)
		boolByte := []byte{0}
		if v {
			boolByte = []byte{1}
		}
		binary.Write(&self.buf, binary.LittleEndian, boolByte)
	case model.FieldType_String:
		self.WriteString(value.Value)
	case model.FieldType_Enum:
		binary.Write(&self.buf, binary.LittleEndian, value.EnumValue)
	default:
		panic("unsupport type" + model.FieldTypeToString(ft))
	}

}

func NewStream() *Stream {
	return &Stream{}
}
