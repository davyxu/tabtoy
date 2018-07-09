package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
)

// 将不同文件/Sheet/KV转换的表，按照表头类型合并数据输出
func mergeData(inputList, outputList *model.DataTableList, symbols *model.TypeTable) {

	for _, inputTab := range inputList.AllTables() {

		var outputTab *model.DataTable

		// 表头类型为分类
		outputTab = outputList.GetDataTable(inputTab.HeaderType)

		// 为输入表头创建唯一的表数据
		if outputTab == nil {
			outputTab = model.NewDataTable()
			outputTab.HeaderType = inputTab.HeaderType
			outputTab.OriginalHeaderType = inputTab.OriginalHeaderType

			// 原始表头类型为解析
			headerFields := symbols.AllFieldByName(inputTab.OriginalHeaderType)

			if headerFields == nil {
				panic("表头找不到!" + inputTab.OriginalHeaderType)
			}

			// 将完整的表头添加到输出表的表头中
			for col, tf := range headerFields {

				// 将输入表表头源信息复制
				inputHeader := inputTab.HeaderByName(tf.Name)

				if inputHeader == nil {
					continue
				}

				outputHeader := outputTab.MustGetHeader(col)
				outputHeader.Cell.Value = tf.Name
				outputHeader.Cell.Col = col
				outputHeader.Cell.Row = 0
				outputHeader.TypeInfo = tf

				headerCell := outputTab.MustGetCell(0, col)
				headerCell.Value = tf.Name

			}

			outputList.AddDataTable(outputTab)
		}

		// 遍历输入表的每一行
		for row := 1; row < len(inputTab.Rows); row++ {

			var outputRow int
			// 输出新开一行

			// 该行有数据，防止注释行加入输出
			if inputTab.GetCell(row, 0) != nil {
				outputRow = outputTab.AddRow()
			}

			// 遍历输入数据的每一列
			for _, inputHeader := range inputTab.Headers {

				// 输入的列头

				if inputHeader.TypeInfo == nil {
					continue
				}

				inputCell := inputTab.GetCell(row, inputHeader.Cell.Col)

				// 这行被注释，无效行
				if inputCell == nil {
					break
				}

				// 用输入的表头名在输出的表头中找
				outputHeader := outputTab.HeaderByName(inputHeader.TypeInfo.FieldName)

				if outputHeader == nil {
					panic("输入的列头名在输出表头中找不到")
				}

				// 取输出表的最后的一行和对应表头的列单元格
				outputCell := outputTab.MustGetCell(outputRow, outputHeader.Cell.Col)

				if outputCell == nil {
					panic("输出单元格找不到")
				}

				if inputHeader.TypeInfo.IsArray() {

					outputCell.Value = combineArrayValue(outputCell.Value, inputCell.Value, inputHeader.TypeInfo.ArraySplitter)

					// TODO 数值来源于多个表格，应该都记录
					outputCell.Table = inputCell.Table
					outputCell.Row = inputCell.Row
					outputCell.Col = inputCell.Col

				} else {
					outputCell.CopyFrom(inputCell)
				}
			}
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
