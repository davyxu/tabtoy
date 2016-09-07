package exportorv2

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

// 从单元格原始数据到最终输出的数值, 检查并转换, 处理默认值及根据meta转换情况
func convertValue(fd *model.FieldDefine, value string, typeset *model.BuildInTypeSet) (ret string, ok bool) {

	// 空格, 且有默认值时, 使用默认值
	if value == "" {
		value = fd.DefaultValue()
	}

	switch fd.Type {
	case model.FieldType_Int32:
		_, err := strconv.ParseInt(value, 10, 32)
		return value, err == nil
	case model.FieldType_Int64:
		_, err := strconv.ParseInt(value, 10, 64)
		return value, err == nil
	case model.FieldType_UInt32:
		_, err := strconv.ParseUint(value, 10, 32)
		return value, err == nil
	case model.FieldType_UInt64:
		_, err := strconv.ParseUint(value, 10, 64)
		return value, err == nil
	case model.FieldType_Float:
		_, err := strconv.ParseFloat(value, 32)
		return value, err == nil
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
				log.Debugln("bool parse failed", err)
				return "", false
			}

			if v {
				ret = "true"
			} else {
				ret = "false"
			}

			break
		}

		return ret, true
	case model.FieldType_String:
		return value, true
	case model.FieldType_Enum:
		if fd.BuildInType == nil {
			log.Errorln("enum type nil", fd.Name)
			return "", false
		}

		evd := fd.BuildInType.FieldByValueAndMeta(value)
		if evd == nil {
			log.Errorf("enum value not found, '%s' enum type: %s", value, fd.BuildInType.Name)
			return "", false
		}

		// 使用枚举的英文字段名输出
		ret = evd.Name

		return ret, true
	case model.FieldType_Struct:
		return "", true

	default:
		log.Errorln("unknown field type in filter", fd.Type, fd.Name)
		return "", false
	}

}
