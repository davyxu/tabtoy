package compiler

import (
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/report"
)

func checkRepeat(globals *model.Globals) {

	for _, tab := range globals.Datas.AllTables() {

		// 遍历输入数据的每一列
		for _, header := range tab.Headers {

			// 输入的列头，为空表示改列被注释
			if header.TypeInfo == nil {
				continue
			}

			// 这列需要建立索引
			if header.TypeInfo.MakeIndex {

				checker := map[string]*model.Cell{}

				for row := 1; row < len(tab.Rows); row++ {

					inputCell := tab.GetCell(row, header.Cell.Col)

					// 这行被注释，无效行
					if inputCell == nil {
						continue
					}

					if _, ok := checker[inputCell.Value]; ok {

						report.ReportError("DuplicateValueInMakingIndex", inputCell.String())

					} else {
						checker[inputCell.Value] = inputCell
					}

				}
			}
		}
	}
}
