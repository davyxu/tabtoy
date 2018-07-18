package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
	"strings"
)

func createOutputTable(symbols *model.TypeTable, inputTab *model.DataTable) *model.DataTable {
	outputTab := model.NewDataTable()
	outputTab.HeaderType = inputTab.HeaderType
	outputTab.OriginalHeaderType = inputTab.OriginalHeaderType

	// 原始表头类型为解析
	headerFields := symbols.AllFieldByName(inputTab.OriginalHeaderType)

	if headerFields == nil {
		report.ReportError("HeaderTypeNotFound", inputTab.OriginalHeaderType)
	}

	// 将完整的表头添加到输出表的表头中
	for col, tf := range headerFields {

		outputHeader := outputTab.MustGetHeader(col)
		outputHeader.Cell.Value = tf.Name
		outputHeader.Cell.Col = col
		outputHeader.Cell.Row = 0
		outputHeader.TypeInfo = tf

		headerCell := outputTab.MustGetCell(0, col)
		headerCell.Value = tf.Name
	}

	return outputTab
}

// 将不同文件/Sheet/KV转换的表，按照表头类型合并数据输出
func MergeData(inputList, outputList *model.DataTableList, symbols *model.TypeTable) {

	for _, inputTab := range inputList.AllTables() {

		var outputTab *model.DataTable

		// 表头类型为分类
		outputTab = outputList.GetDataTable(inputTab.HeaderType)

		// 为输入表头创建唯一的表数据
		if outputTab == nil {
			outputTab = createOutputTable(symbols, inputTab)

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
					panic("输入的列头名在输出表头中找不到:" + inputHeader.TypeInfo.FieldName)
				}

				// 取输出表的最后的一行和对应表头的列单元格
				outputCell := outputTab.MustGetCell(outputRow, outputHeader.Cell.Col)

				if outputCell == nil {
					panic("输出单元格找不到")
				}

				if inputHeader.TypeInfo.IsArray() {

					combineArrayCell(outputCell, inputCell, inputHeader.TypeInfo.ArraySplitter)

				} else {
					outputCell.CopyFrom(inputCell)
				}
			}
		}

	}
}

func combineArrayCell(ouputCell, inputCell *model.Cell, splitter string) {

	var sb strings.Builder

	var valueCount int

	tail := ouputCell

	// 把之前的格子的值合并为字符串
	for c := ouputCell.Next; c != nil; c = c.Next {

		tail = c

		if valueCount > 0 {
			sb.WriteString(splitter)
		}

		sb.WriteString(c.Value)
		valueCount++
	}

	for _, value := range strings.Split(inputCell.Value, splitter) {

		if valueCount > 0 {
			sb.WriteString(splitter)
		}

		sb.WriteString(value)
		valueCount++
	}

	ouputCell.Value = sb.String()

	tail.Next = &model.Cell{}
	tail.Next.CopyFrom(inputCell)

}
