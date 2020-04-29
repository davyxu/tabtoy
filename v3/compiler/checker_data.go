package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"strconv"
	"strings"
)

// 检查数据与定义类型是否匹配
func checkDataType(globals *model.Globals) {

	var currHeader *model.HeaderField
	var crrCell *model.Cell

	for _, tab := range globals.Datas.AllTables() {

		// 遍历输入数据的每一列
		for _, header := range tab.Headers {

			// 输入的列头，为空表示改列被注释
			if header.TypeInfo == nil {
				continue
			}

			for row := 1; row < len(tab.Rows); row++ {

				inputCell := tab.GetCell(row, header.Cell.Col)

				// 这行被注释，无效行
				if inputCell == nil {
					continue
				}

				crrCell = inputCell
				currHeader = header

				if inputCell.Value == "" {
					continue
				}

				if header.TypeInfo.IsArray() {
					for _, value := range strings.Split(inputCell.Value, header.TypeInfo.ArraySplitter) {
						err := checkSingleValue(header, value)
						if err != nil {
							report.ReportError("DataMissMatchTypeDefine", currHeader.TypeInfo.FieldType, crrCell.String())
						}
					}
				} else {
					err := checkSingleValue(header, inputCell.Value)
					if err != nil {
						report.ReportError("DataMissMatchTypeDefine", currHeader.TypeInfo.FieldType, crrCell.String())
					}
				}

			}
		}
	}
}

func checkSingleValue(header *model.HeaderField, value string) error {
	switch model.LanguagePrimitive(header.TypeInfo.FieldType, "go") {
	case "int16":
		_, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return err
		}
	case "int32":
		_, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
	case "int64":
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
	case "uint16":
		_, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return err
		}
	case "uint32":
		_, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
	case "uint64":
		_, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
	case "float32":
		_, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
	case "float64":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
	case "bool":
		_, err := model.ParseBool(value)
		if err != nil {
			return err
		}
	}

	return nil
}
