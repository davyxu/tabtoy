package bindata

import (
	"bytes"
	"encoding/binary"
	"math"
)

type BinaryWriter struct {
	buffer bytes.Buffer
}

func (self *BinaryWriter) Bytes() []byte {
	return self.buffer.Bytes()
}

func (self *BinaryWriter) Write(b []byte) (int, error) {
	return self.buffer.Write(b)
}

func (self *BinaryWriter) WriteInt16(x int16) error {

	return self.WriteUInt16(uint16(x))
}

func (self *BinaryWriter) WriteInt32(x int32) error {
	return self.WriteUInt32(uint32(x))
}

func (self *BinaryWriter) WriteInt64(x int64) error {
	return self.WriteUInt64(uint64(x))
}

func (self *BinaryWriter) WriteUInt16(x uint16) error {

	return binary.Write(&self.buffer, binary.LittleEndian, &x)
}

func (self *BinaryWriter) WriteUInt32(x uint32) error {

	return binary.Write(&self.buffer, binary.LittleEndian, &x)
}

func (self *BinaryWriter) WriteUInt64(x uint64) error {
	return binary.Write(&self.buffer, binary.LittleEndian, &x)
}

func (self *BinaryWriter) WriteFloat32(x float32) error {

	v := math.Float32bits(x)
	return binary.Write(&self.buffer, binary.LittleEndian, &v)
}

func (self *BinaryWriter) WriteFloat64(x float64) error {

	v := math.Float64bits(x)
	return binary.Write(&self.buffer, binary.LittleEndian, &v)
}

func (self *BinaryWriter) WriteBool(x bool) error {
	return binary.Write(&self.buffer, binary.LittleEndian, &x)
}

func (self *BinaryWriter) WriteString(x string) error {
	l := uint32(len(x))
	if err := binary.Write(&self.buffer, binary.LittleEndian, &l); err != nil {
		return err
	}

	data := []byte(x)
	if err := binary.Write(&self.buffer, binary.LittleEndian, &data); err != nil {
		return err
	}

	return nil
}

func NewBinaryWriter() *BinaryWriter {
	return &BinaryWriter{}
}
