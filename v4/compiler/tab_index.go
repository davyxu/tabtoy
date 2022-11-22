package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"strings"
)

const (
	idxHeaderCol_Mode = 0
	idxHeaderCol_Type = 1
	idxHeaderCol_File = 2
	maxIndexHeaderCol = 3
)

func loadIndexHeader(sheet util.TableSheet) (colByHeaderType [maxIndexHeaderCol]int, ok bool) {
	for col := 0; col < maxIndexHeaderCol; col++ {

		headerValue := sheet.GetValue(0, col, nil)

		var headerType int
		switch headerValue {
		case "Mode":
			headerType = idxHeaderCol_Mode
		case "Type":
			headerType = idxHeaderCol_Type
		case "File":
			headerType = idxHeaderCol_File
		default:
			return
		}

		colByHeaderType[headerType] = col
	}

	ok = true

	return
}

func loadIndexTable(file util.TableFile, fileName string) (ret []*model.FileMeta) {
	for _, sheet := range file.Sheets() {

		colByHeaderType, ok := loadIndexHeader(sheet)

		if !ok {
			util.ReportError("InvalidIndexHeader", fileName)
			return
		}

		// 遍历所有数据行
		for row := 1; ; row++ {
			if sheet.IsRowEmpty(row, maxIndexHeaderCol+1) {
				break
			}

			firstCol := sheet.GetValue(row, 0, nil)
			// 首列带#时，本行忽略
			if strings.HasPrefix(firstCol, "#") {
				continue
			}

			var meta model.FileMeta
			meta.FileName = sheet.GetValue(row, colByHeaderType[idxHeaderCol_File], nil)
			meta.Mode = sheet.GetValue(row, colByHeaderType[idxHeaderCol_Mode], nil)
			meta.HeaderType = sheet.GetValue(row, colByHeaderType[idxHeaderCol_Type], nil)

			switch meta.Mode {
			case "Data", "KV":
				if meta.HeaderType == "" {
					util.ReportError("EmptyTableType", fileName)
					return
				}

			case "Type":

			default:
				util.ReportError("InvalidTableMode", fileName)
				return
			}

			ret = append(ret, &meta)

		}
	}

	return
}
