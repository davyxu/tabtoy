package printer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
)

type BinaryFile struct {
	Name string // rootName
	buf  bytes.Buffer
}

func (self *BinaryFile) Data() []byte {
	return self.buf.Bytes()
}

func (self *BinaryFile) Buffer() *bytes.Buffer {
	return &self.buf
}

func (self *BinaryFile) Printf(format string, args ...interface{}) {
	self.buf.WriteString(fmt.Sprintf(format, args...))
}

func (self *BinaryFile) Write(outfile string) bool {
	err := ioutil.WriteFile(outfile, self.buf.Bytes(), 0666)
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	return true
}

func (self *BinaryFile) WriteInt32(v int32) {
	binary.Write(&self.buf, binary.LittleEndian, v)
}

func (self *BinaryFile) WriteString(v string) {
	rawStr := []byte(v)

	binary.Write(&self.buf, binary.LittleEndian, int32(len(rawStr)))
	binary.Write(&self.buf, binary.LittleEndian, rawStr)
}

func NewBinaryFile(name string) *BinaryFile {
	return &BinaryFile{
		Name: name,
	}
}
