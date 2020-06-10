package jsondata2

import (
	"encoding/json"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"strconv"
	"strings"
)

func wrapValue(globals *model.Globals, value string, valueType *model.TypeDefine) interface{} {
	if valueType.IsArray() {

		var vlist = make([]interface{}, 0)
		// 空的单元格，导出空数组，除非强制指定填充默认值
		if value != "" {
			for _, elementValue := range strings.Split(value, valueType.ArraySplitter) {

				vlist = append(vlist, wrapSingleValue(globals, valueType, elementValue))
			}
		}

		return vlist

	} else {
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

func Generate(globals *model.Globals) (data []byte, err error) {

	fileData := map[string]interface{}{
		"@Tool":    "github.com/davyxu/tabtoy",
		"@Version": globals.Version,
	}

	for _, tab := range globals.Datas.AllTables() {

		headers := globals.Types.AllFieldByName(tab.OriginalHeaderType)

		var tabData []interface{}

		for row := 1; row < len(tab.Rows); row++ {

			rowData := map[string]interface{}{}
			for col, header := range headers {

				// 在单元格找到值
				valueCell := tab.GetCell(row, col)

				var value interface{}
				if valueCell != nil {

					value = wrapValue(globals, valueCell.Value, header)
				} else {
					// 这个表中没有这列数据
					value = wrapValue(globals, "", header)
				}

				rowData[header.FieldName] = value
			}

			tabData = append(tabData, rowData)

		}

		fileData[tab.HeaderType] = tabData
	}

	return json.MarshalIndent(&fileData, "", "\t")
}
