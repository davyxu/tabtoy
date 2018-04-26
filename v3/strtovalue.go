package v3

import (
	"reflect"
	"strconv"
)

func StringToValue(str string, value interface{}) error {

	switch raw := value.(type) {
	case *int32:
		v, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err
		}

		*raw = int32(v)
	case *int64:
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return err
		}

		*raw = v
	case *uint32:
		v, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			return err
		}

		*raw = uint32(v)
	case *uint64:
		v, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return err
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
				return err
			}
		}

		*raw = v
	case *float32:
		v, err := strconv.ParseFloat(str, 32)
		if err != nil {
			return err
		}

		*raw = float32(v)
	case *float64:
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return err
		}

		*raw = float64(v)

	default:
		panic("unsupport type: " + reflect.TypeOf(value).Elem().Name())
	}

	return nil
}
