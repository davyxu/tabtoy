package jsondata

import (
	"encoding/json"
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
)

func OutputDir(globals *model.Globals, outputDir string) (err error) {

	for _, tab := range globals.Datas.AllTables() {

		// 一个表的所有列
		headers := globals.Types.AllFieldByName(tab.HeaderType)

		fileData := map[string]interface{}{
			"@Tool":    "github.com/davyxu/tabtoy",
			"@Version": globals.Version,
		}

		var tabData []interface{}

		// 遍历所有行
		for row := 0; row < len(tab.Rows); row++ {

			// 遍历每一列
			rowData := map[string]interface{}{}
			for col, header := range headers {

				// 在单元格找到值
				valueCell := tab.GetCell(row, col)

				var value = wrapValue(globals, valueCell, header)

				rowData[header.FieldName] = value
			}

			tabData = append(tabData, rowData)
		}

		fileData[tab.HeaderType] = tabData

		var data []byte
		data, err = json.MarshalIndent(&fileData, "", "\t")

		if err != nil {
			return err
		}

		err = util.WriteFile(fmt.Sprintf("%s/%s.json", outputDir, tab.HeaderType), data)
		if err != nil {
			return err
		}
	}

	return nil
}
func OutputFile(globals *model.Globals, outFile string) (err error) {

	data, err := OutputData(globals)
	if err != nil {
		return err
	}

	return util.WriteFile(outFile, data)
}

func OutputData(globals *model.Globals) (data []byte, err error) {

	fileData := map[string]interface{}{
		"@Tool":    "github.com/davyxu/tabtoy",
		"@Version": globals.Version,
	}

	for _, tab := range globals.Datas.AllTables() {

		// 一个表的所有列
		headers := globals.Types.AllFieldByName(tab.HeaderType)

		var tabData []interface{}

		// 遍历所有行
		for row := 0; row < len(tab.Rows); row++ {

			// 遍历每一列
			rowData := map[string]interface{}{}
			for col, header := range headers {

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
