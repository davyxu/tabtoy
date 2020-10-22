package jsondata

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
)

func Output(globals *model.Globals, param string) (err error) {

	for _, tab := range globals.Datas.AllTables() {

		// 一个表的所有列
		headers := globals.Types.AllFieldByName(tab.OriginalHeaderType)

		fileData := map[string]interface{}{
			"@Tool":    "github.com/davyxu/tabtoy",
			"@Version": globals.Version,
		}

		var tabData []interface{}

		// 遍历所有行
		for row := 1; row < len(tab.Rows); row++ {

			// 遍历每一列
			rowData := map[string]interface{}{}
			for col, header := range headers {

				if globals.CanDoAction(model.ActionNoGenJsonDir, header) {
					continue
				}

				// 在单元格找到值
				valueCell := tab.GetCell(row, col)

				var value = wrapValue(globals, valueCell, header)

				rowData[header.FieldName] = value
			}

			tabData = append(tabData, rowData)
		}

		fileData[tab.HeaderType] = tabData

		data, err := json.MarshalIndent(&fileData, "", "\t")

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", param, tab.HeaderType), data, 0666)
		if err != nil {
			return err
		}
	}

	return nil
}

func Generate(globals *model.Globals) (data []byte, err error) {

	fileData := map[string]interface{}{
		"@Tool":    "github.com/davyxu/tabtoy",
		"@Version": globals.Version,
	}

	for _, tab := range globals.Datas.AllTables() {

		// 一个表的所有列
		headers := globals.Types.AllFieldByName(tab.OriginalHeaderType)

		var tabData []interface{}

		// 遍历所有行
		for row := 1; row < len(tab.Rows); row++ {

			// 遍历每一列
			rowData := map[string]interface{}{}
			for col, header := range headers {

				if globals.CanDoAction(model.ActionNoGenJson, header) {
					continue
				}

				// 在单元格找到值
				valueCell := tab.GetCell(row, col)

				var value = wrapValue(globals, valueCell, header)

				rowData[header.FieldName] = value
			}

			tabData = append(tabData, rowData)
		}

		fileData[tab.HeaderType] = tabData
	}

	return json.MarshalIndent(fileData, "", "\t")
}
