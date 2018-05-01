package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
	"strings"
)

// 将不同文件/Sheet/KV转换的表，按照表头类型合并数据输出
func mergeData(inputList, outputList *model.DataTableList, symbols *model.SymbolTable) {

	for _, inputTab := range inputList.Data {

		var outputTab *model.DataTable

		// 表头类型为分类
		outputTab = outputList.GetDataTable(inputTab.HeaderType)

		// 为输入表头创建唯一的表数据
		if outputTab == nil {
			outputTab = model.NewDataTable()
			outputTab.HeaderType = inputTab.HeaderType
			outputTab.OriginalHeaderType = inputTab.OriginalHeaderType

			var headerFields []*table.TableField

			// 原始表头类型为解析
			headerFields = symbols.Fields(inputTab.OriginalHeaderType)

			if headerFields == nil {
				panic("表头找不到!" + inputTab.OriginalHeaderType)
			}

			outputTab.HeaderFields = headerFields

			outputList.AddDataTable(outputTab)
		}

		// 遍历输入表的每一行
		for _, rowData := range inputTab.Rows {

			// 提前准备完整列头的宽度
			oneRow := make(model.DataRow, len(outputTab.HeaderFields))

			// 遍历输入数据的每一列
			for col, value := range rowData {

				// 输入的列头
				headerField := inputTab.HeaderFields[col]

				// 用输入的表头名在输出的表头中找
				_, OutputCol := outputTab.HeaderFieldByName(headerField.FieldName)

				if headerField.IsArray() {

					oneRow[OutputCol] = combineArrayValue(oneRow[OutputCol], value, headerField.ArraySplitter)

				} else {
					oneRow[OutputCol] = value
				}
			}

			outputTab.AddRow(oneRow)
		}

	}
}

func combineArrayValue(oldValue, newValue, splitter string) string {
	var oldvalues, newvalues []string

	if oldValue != "" {
		oldvalues = strings.Split(oldValue, splitter)
	}

	if newValue != "" {
		newvalues = strings.Split(newValue, splitter)
	}

	var sb strings.Builder
	for index, str := range append(oldvalues, newvalues...) {
		if index > 0 {
			sb.WriteString(splitter)
		}
		sb.WriteString(str)
	}

	return sb.String()
}
