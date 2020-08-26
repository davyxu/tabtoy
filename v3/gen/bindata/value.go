package bindata

import (
	"github.com/davyxu/tabtoy/v3/model"
	"strconv"
)

func writeValue(globals *model.Globals, structWriter *BinaryWriter, fieldType *model.TypeDefine, goType, value string) error {
	switch {
	case goType == "int16":
		if value == "" {
			return structWriter.WriteInt16(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}

			return structWriter.WriteInt16(int16(v))
		}

	case goType == "int32":
		if value == "" {
			return structWriter.WriteInt32(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}

			return structWriter.WriteInt32(int32(v))
		}

	case goType == "int64":
		if value == "" {
			return structWriter.WriteInt64(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}

			return structWriter.WriteInt64(int64(v))
		}

	case goType == "uint16":
		if value == "" {
			return structWriter.WriteUInt16(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return err
			}

			return structWriter.WriteUInt16(uint16(v))
		}

	case goType == "uint32":
		if value == "" {
			return structWriter.WriteUInt32(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return err
			}

			return structWriter.WriteUInt32(uint32(v))
		}

	case goType == "uint64":
		if value == "" {
			return structWriter.WriteUInt64(0)
		} else {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}

			return structWriter.WriteUInt64(uint64(v))
		}

	case goType == "bool":
		if value == "" {
			return structWriter.WriteBool(false)
		} else {
			var v bool
			switch value {
			case "是", "yes", "YES", "1", "true", "TRUE", "True":
				v = true
			default:
			}

			return structWriter.WriteBool(v)
		}

	case goType == "float32":
		if value == "" {
			return structWriter.WriteFloat32(0)
		} else {
			v, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return err
			}

			return structWriter.WriteFloat32(float32(v))
		}
	case goType == "float64":
		if value == "" {
			return structWriter.WriteFloat64(0)
		} else {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			}

			return structWriter.WriteFloat64(v)
		}
	case goType == "string":
		return structWriter.WriteString(value)

	case globals.Types.IsEnumKind(fieldType.FieldType): // 枚举
		if value == "" {
			return structWriter.WriteInt32(0)
		} else {
			enumValue := globals.Types.ResolveEnumValue(fieldType.FieldType, value)

			v, err := strconv.ParseInt(enumValue, 10, 32)
			if err != nil {
				return err
			}

			return structWriter.WriteInt32(int32(v))
		}

	default:
		panic("unknown binary type: " + fieldType.FieldType)
	}

	return nil
}
