package golang

import "encoding/binary"

type binaryReader struct {
	buf []byte
}

func (self *binaryReader) ReadBool(x *bool) {
	n := self.buf[0]
	self.buf = self.buf[1:]
	*x = n != 0
}

func (self *binaryReader) ReadUInt16(x *uint16) {
	*x = binary.LittleEndian.Uint16(self.buf[0:2])
	self.buf = self.buf[2:]
}

func (self *binaryReader) ReadUInt32(x *uint32) {
	*x = binary.LittleEndian.Uint32(self.buf[0:4])
	self.buf = self.buf[4:]
}
func (self *binaryReader) ReadUInt64(x *uint64) {
	*x = binary.LittleEndian.Uint64(self.buf[0:8])
	self.buf = self.buf[8:]
}

func (self *binaryReader) ReadBytes(x *[]byte) {
	var l uint16
	self.ReadUInt16(&l)
	*x = make([]byte, l)

	copy(*x, self.buf[0:l])
	self.buf = self.buf[l:]
}

func (self *binaryReader) ReadInt16(x *int16) {

	var v uint16
	self.ReadUInt16(&v)

	*x = int16(v)
}

func (self *binaryReader) ReadInt32(x *int32) {
	var v uint32
	self.ReadUInt32(&v)

	*x = int32(v)
}

func (self *binaryReader) ReadInt64(x *int64) {
	var v uint64
	self.ReadUInt64(&v)

	*x = int64(v)
}
