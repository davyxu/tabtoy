package model

func IsNativeType(typeName string) bool {

	switch typeName {
	case "int32", "uint32", "int64", "uint64", "float", "bool", "string":
		return true
	}

	return false
}
