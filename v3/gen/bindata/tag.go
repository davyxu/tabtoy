package bindata

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func MakeTag(globals *model.Globals, tf *model.TypeDefine, fieldIndex int) uint32 {
	convertedType := model.LanguagePrimitive(tf.FieldType, "go")

	var t int
	switch {
	case convertedType == "int16":
		t = 1
	case convertedType == "int32":
		t = 2
	case convertedType == "int64":
		t = 3
	case convertedType == "uint16":
		t = 4
	case convertedType == "uint32":
		t = 5
	case convertedType == "uint64":
		t = 6
	case convertedType == "float32":
		t = 7
	case convertedType == "string":
		t = 8
	case convertedType == "bool":
		t = 9
	case globals.Types.IsEnumKind(tf.FieldType):
		t = 10
	case convertedType == "float64":
		t = 12
		// 注意, t = 11是结构体
	default:
		panic("unknown type:" + tf.FieldType)
	}

	if tf.IsArray() {
		t += 100
	}

	return uint32(t<<16 | fieldIndex)
}

func MakeTagStructArray() uint32 {

	var t int
	t = 11

	// 结构体默认是数组
	t += 100

	return uint32(t << 16)
}

func writePair(globals *model.Globals, structWriter *BinaryWriter, fieldType *model.TypeDefine, goType, value string, fieldIndex int) error {

	tag := MakeTag(globals, fieldType, fieldIndex)
	if err := structWriter.WriteUInt32(tag); err != nil {
		return err
	}

	return writeValue(globals, structWriter, fieldType, goType, value)
}
