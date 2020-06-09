package filter

import (
	"strconv"

	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

// 从单元格原始数据到最终输出的数值, 检查并转换, 处理默认值及根据meta转换情况
func ConvertValue(fd *model.FieldDescriptor, value string, fileD *model.FileDescriptor, node *model.Node) (ret string, ok bool) {

	// 空格, 且有默认值时, 使用默认值
	if value == "" {
		value = fd.DefaultValue()
	}

	switch fd.Type {
	case model.FieldType_Int32:
		v, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_Int64:
		v, err := strconv.ParseInt(value, 10, 64)

		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_UInt32:
		v, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_UInt64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_Float:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_Bool:
		if value == "是" || value == "true" || value == "True" {
			value = "true"
		} else {
			value = "false"
		}

		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Debugln(err)
			return "", false
		}

		ret = value
		node.AddValue(ret, v)
	case model.FieldType_String:
		ret = value
		node.AddValue(ret, value)
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
		node.AddValue(ret, evd.EnumValue).EnumValue = evd.EnumValue

	case model.FieldType_Struct:

		if fd.Complex == nil {
			log.Errorf("%s, '%s'", i18n.String(i18n.ConvertValue_StructTypeNil), fd.Name)
			return "", false
		}

		if value == "" {

			if !fillStructDefaultValue(fd.Complex, fileD, node) {
				return "", false
			}

		} else {
			if !parseStruct(fd, value, fileD, node) {
				return "", false
			}
		}

	default:
		log.Errorf("%s, '%s' '%s'", i18n.String(i18n.ConvertValue_UnknownFieldType), fd.Name, fd.Name)
		return "", false
	}

	ok = true

	return
}

// 填充空结构体的默认值
func fillStructDefaultValue(structD *model.Descriptor, fileD *model.FileDescriptor, node *model.Node) bool {

	for _, fd := range structD.Fields {

		// 没默认值不输出, 建议忽略的字段除外, 先导出node, 再在printer中忽略
		if fd.Meta.GetString("Default") == "" && node.SugguestIgnore {
			continue
		}

		fieldNode := node.AddKey(fd)

		// 结构体的值没填, 且没默认值, 建议忽略
		if fd.Meta.GetString("Default") == "" && node.Value == "" {
			fieldNode.SugguestIgnore = true
		}

		_, ok := ConvertValue(fd, "", fileD, fieldNode)
		if !ok {
			return false
		}
	}

	return true

}
