package jsondata

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"strconv"
)

func wrapValue(globals *model.Globals, valueCell *model.Cell, valueType *model.TypeDefine) interface{} {
	if valueType.IsArray() {

		var vlist = make([]interface{}, 0)
		// 空的单元格，导出空数组，除非强制指定填充默认值
		if valueCell != nil {

			for _, elementValue := range valueCell.ValueList {

				vlist = append(vlist, wrapSingleValue(globals, valueType, elementValue))
			}
		}

		return vlist

	} else {

		var value string
		if valueCell != nil {
			value = valueCell.Value
		}

		return wrapSingleValue(globals, valueType, value)
	}
}

func wrapSingleValue(globals *model.Globals, valueType *model.TypeDefine, value string) interface{} {

	goType := model.LanguagePrimitive(valueType.FieldType, "go")

	switch {
	case goType == "string": // 字符串
		return util.StringEscape(value)
	case goType == "float32":

		if value == "" {
			return float32(0)
		}

		f64, _ := strconv.ParseFloat(value, 32)
		return float32(f64)
	case goType == "float64":

		if value == "" {
			return float64(0)
		}

		f64, _ := strconv.ParseFloat(value, 64)
		return f64
	case globals.Types.IsEnumKind(valueType.FieldType): // 枚举
		enumValue := globals.Types.ResolveEnumValue(valueType.FieldType, value)
		i, _ := strconv.Atoi(enumValue)
		return int32(i)
	case goType == "bool":

		v, _ := model.ParseBool(value)
		if v {
			return true
		}

		return false
	case goType == "int16":
		i64, _ := strconv.ParseInt(value, 10, 16)
		return int16(i64)
	case goType == "int32":
		i64, _ := strconv.ParseInt(value, 10, 32)
		return int32(i64)
	case goType == "int64":
		i64, _ := strconv.ParseInt(value, 10, 64)
		return i64
	case goType == "uint16":
		i64, _ := strconv.ParseInt(value, 10, 16)
		return uint16(i64)
	case goType == "uint32":
		i64, _ := strconv.ParseUint(value, 10, 32)
		return uint32(i64)
	case goType == "uint64":
		i64, _ := strconv.ParseUint(value, 10, 64)
		return i64
	}

	if value == "" {
		return model.FetchDefaultValue(valueType.FieldType)
	}

	return value
}
