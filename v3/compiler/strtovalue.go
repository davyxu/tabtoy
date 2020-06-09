package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"reflect"
	"strconv"
	"strings"
)

func StringToValue(str string, value interface{}, tf *model.TypeDefine, symbols *model.TypeTable) error {

	err, handled := util.StringToPrimitive(str, value)
	if err != nil || handled {
		return err
	}

	if tf == nil {
		panic("unsupport type: " + reflect.TypeOf(value).Elem().Name())
	}

	if tf.IsArray() {

		// Tags为空时, 数组中就不应该有元素
		if str == "" {
			return nil
		}

		tValue := reflect.TypeOf(value).Elem()
		vValue := reflect.Indirect(reflect.ValueOf(value))

		if vValue.Kind() != reflect.Slice {
			panic("require slice" + str)
		}

		splitedData := strings.Split(str, tf.ArraySplitter)

		slice := reflect.MakeSlice(tValue, len(splitedData), len(splitedData))

		for index, strValue := range splitedData {

			elemElem := slice.Index(index)
			err, handled = util.StringToPrimitive(strValue, elemElem.Addr().Interface())
			if err != nil {
				return err
			}

		}

		vValue.Set(slice)

		return nil
	}

	if symbols.IsEnumKind(tf.FieldType) {

		enumValueStr := symbols.ResolveEnumValue(tf.FieldType, str)

		if enumValueStr != "" {
			enumValue, err := strconv.Atoi(enumValueStr)
			if err != nil {
				return err
			}
			vValue := reflect.Indirect(reflect.ValueOf(value))
			vValue.SetInt(int64(enumValue))
		}

		return nil
	}

	// 内建类型定义与model.InitBuiltinTypes中不匹配
	panic("unhandled value: " + str)

	return nil
}
