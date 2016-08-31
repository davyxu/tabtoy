package filter

import (
	"strconv"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
)

func ValueConvetor(fd *pbmeta.FieldDescriptor, value string) (string, bool) {

	// 空单元格时， 给定一个这个类型对应的值
	if value == "" {
		value = data.GetDefaultValue(fd)
	}

	switch fd.Type() {
	case pbprotos.FieldDescriptorProto_TYPE_FLOAT:

		_, err := strconv.ParseFloat(value, 32)

		if err != nil {

			return "", false
		}

	case pbprotos.FieldDescriptorProto_TYPE_INT64:

		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {

			return "", false
		}

	case pbprotos.FieldDescriptorProto_TYPE_UINT64:

		_, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return "", false
		}

	case pbprotos.FieldDescriptorProto_TYPE_INT32:

		_, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return "", false
		}

	case pbprotos.FieldDescriptorProto_TYPE_UINT32:

		_, err := strconv.ParseUint(value, 10, 32)

		if err != nil {
			return "", false
		}

	case pbprotos.FieldDescriptorProto_TYPE_BOOL:

		var final string

		for {
			if value == "是" {
				final = "true"
				break
			} else if value == "否" {
				final = "false"
				break
			}

			v, err := strconv.ParseBool(value)

			if err != nil {
				return "", false
			}

			if v {
				final = "true"
			} else {
				final = "false"
			}

			break
		}

		value = final

	case pbprotos.FieldDescriptorProto_TYPE_ENUM:
		ed := fd.EnumDesc()

		if ed.ValueCount() == 0 {
			return "", false
		}

		// 枚举值从表格读出时可能是中文枚举值， 需要根据meta信息转换为程序枚举值
		var convValue string = value

		// 遍历这个枚举类型
		for i := 0; i < ed.ValueCount(); i++ {

			evd := ed.Value(i)

			// 取出每个字段的meta
			meta, ok := data.GetFieldMeta(evd)

			if !ok {
				return "", false
			}

			if meta == nil {
				continue
			}

			// 这个枚举值的别名是否与给定的一致
			if meta.Alias == value {
				convValue = evd.Name()
				break
			}
		}

		if ed.ValueByName(convValue) == nil {
			log.Errorf("enum doesn't contain this value, %s in %s", convValue, fd.Name())
			return "", false
		}

		value = convValue

	}

	return value, true
}
