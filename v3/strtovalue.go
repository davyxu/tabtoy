package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
	"strconv"
	"strings"
)

func RawStringToValue(str string, value interface{}) (error, bool) {
	switch raw := value.(type) {
	case *int32:
		v, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err, false
		}

		*raw = int32(v)
	case *int64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err, false
		}

		*raw = v
	case *uint32:
		v, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return err, false
		}

		*raw = uint32(v)
	case *uint64:
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err, false
		}

		*raw = v
	case *string:
		*raw = str
	case *bool:

		var v bool
		var err error

		switch str {
		case "是":
			v = true
		case "否":
			v = false
		default:
			v, err = strconv.ParseBool(str)
			if err != nil {
				return err, false
			}
		}

		*raw = v
	case *float32:
		v, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err, false
		}

		*raw = float32(v)
	case *float64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err, false
		}

		*raw = float64(v)

	default:
		return nil, false
	}

	return nil, true
}

func StringToValue(str string, value interface{}, tf *table.TableField, symbols *model.SymbolTable) error {

	err, handled := RawStringToValue(str, value)
	if err != nil || handled {
		return err
	}

	if tf == nil {
		panic("unsupport type: " + reflect.TypeOf(value).Elem().Name())
	}

	if tf.IsArray && tf.Splitter != "" {

		tValue := reflect.TypeOf(value).Elem()
		vValue := reflect.Indirect(reflect.ValueOf(value))

		if vValue.Kind() != reflect.Slice {
			panic("require slice" + str)
		}

		splitedData := strings.Split(str, tf.Splitter)

		slice := reflect.MakeSlice(tValue, len(splitedData), len(splitedData))

		for index, strValue := range splitedData {

			elemElem := slice.Index(index)
			err, handled = RawStringToValue(strValue, elemElem.Addr().Interface())
			if err != nil {
				return err
			}

		}

		vValue.Set(slice)

		return nil
	}

	if symbols.IsEnumKind(tf.FieldType) {

		enumValue, err := strconv.Atoi(symbols.ResolveEnumValue(tf.FieldType, str))
		if err != nil {
			return err
		}
		vValue := reflect.Indirect(reflect.ValueOf(value))
		vValue.SetInt(int64(enumValue))

		return nil
	}

	panic("unhandled value: " + str)

	return nil
}
