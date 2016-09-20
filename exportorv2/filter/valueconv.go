package filter

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

// 从单元格原始数据到最终输出的数值, 检查并转换, 处理默认值及根据meta转换情况
func ConvertValue(fd *model.FieldDescriptor, value string, fileD *model.FileDescriptor, node *model.Node) (ret string, ok bool) {

	// 空格, 且有默认值时, 使用默认值
	if value == "" {
		value = fd.DefaultValue()
	}

	switch fd.Type {
	case model.FieldType_Int32:
		_, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret)
	case model.FieldType_Int64:
		_, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret)
	case model.FieldType_UInt32:
		_, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret)
	case model.FieldType_UInt64:
		_, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret)
	case model.FieldType_Float:
		_, err := strconv.ParseFloat(value, 32)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret)
	case model.FieldType_Bool:

		for {
			if value == "是" {
				ret = "true"
				break
			} else if value == "否" {
				ret = "false"
				break
			}

			v, err := strconv.ParseBool(value)

			if err != nil {
				log.Debugln(err)
				return "", false
			}

			if v {
				ret = "true"
			} else {
				ret = "false"
			}

			break
		}

		node.AddValue(ret)

	case model.FieldType_String:
		ret = value
		node.AddValue(ret)
	case model.FieldType_Enum:
		if fd.Complex == nil {
			log.Errorf("%s, '%s'", i18n.String(i18n.ConvertValue_EnumTypeNil), fd.Name)
			return "", false
		}

		evd := fd.Complex.FieldByValueAndMeta(value)
		if evd == nil {
			log.Errorf("%s, '%s' '%s'", i18n.String(i18n.ConvertValue_EnumValueNotFound), value, fd.Complex.Name)
			return "", false
		}

		// 使用枚举的英文字段名输出
		ret = evd.Name
		node.AddValue(ret).EnumValue = evd.EnumValue

	case model.FieldType_Struct:

		if fd.Complex == nil {
			log.Errorf("%s, '%s'", i18n.String(i18n.ConvertValue_StructTypeNil), fd.Name)
			return "", false
		}

		if !parseStruct(fd, value, fileD, node) {
			return "", false
		}

	default:
		log.Errorf("%s, '%s' '%s'", i18n.String(i18n.ConvertValue_UnknownFieldType), fd.Name, fd.Name)
		return "", false
	}

	ok = true

	return
}
