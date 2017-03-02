package printer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type BinaryFile struct {
	buf       bytes.Buffer
	stopAtPos int
}

func (self *BinaryFile) Len() int {
	return self.buf.Len()
}

func (self *BinaryFile) Buffer() *bytes.Buffer {
	return &self.buf
}

func (self *BinaryFile) WriteBytes(b []byte) {
	self.buf.Write(b)
	self.checkStopPos()
}

func (self *BinaryFile) Printf(format string, args ...interface{}) {
	self.buf.WriteString(fmt.Sprintf(format, args...))
	self.checkStopPos()
}

func (self *BinaryFile) Write(outfile string) bool {
	err := ioutil.WriteFile(outfile, self.buf.Bytes(), 0666)
	if err != nil {
		log.Errorf("%s, %v", i18n.String(i18n.Printer_OpenWriteOutputFileFailed), err.Error())
		return false
	}

	return true
}

func (self *BinaryFile) checkStopPos() {
	if self.stopAtPos != -1 && self.Len() == self.stopAtPos {
		self.stopAtPos = self.stopAtPos
	}
}

func (self *BinaryFile) StopAtPos(pos int) {
	self.stopAtPos = pos
}

func (self *BinaryFile) WriteInt32(v int32) {

	binary.Write(&self.buf, binary.LittleEndian, v)
	self.checkStopPos()
}

func (self *BinaryFile) WriteString(v string) {
	rawStr := []byte(v)

	binary.Write(&self.buf, binary.LittleEndian, int32(len(rawStr)))

	self.checkStopPos()

	binary.Write(&self.buf, binary.LittleEndian, rawStr)

	self.checkStopPos()
}

func (self *BinaryFile) WriteNodeValue(ft model.FieldType, value *model.Node) {

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

	self.checkStopPos()

}

func NewBinaryFile() *BinaryFile {
	return &BinaryFile{
		stopAtPos: -1,
	}
}
