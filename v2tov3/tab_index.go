package v2tov3

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"sort"
)

func ExportIndexTable(globals *model.Globals) error {

	globals.TargetIndexSheet = globals.AddTable("Index.xlsx")

	helper.WriteIndexTableHeader(globals.TargetIndexSheet)

	var tabList []*helper.MemFileData

	globals.TargetTables.VisitAllTable(func(data *helper.MemFileData) bool {

		if data.FileName == "Index.xlsx" {
			return true
		}

		tabList = append(tabList, data)

		return true
	})

	// 内容排序
	sort.SliceStable(tabList, func(i, j int) bool {

		a := tabList[i]
		b := tabList[j]

		if aMode, bMode := getMode(a), getMode(b); aMode != bMode {

			if aMode == "类型表" {
				return true
			}

			if aMode == "数据表" {
				return false
			}

		}

		if a.TableName != b.TableName {
			return a.TableName < b.TableName
		}

		return a.FileName < b.FileName
	})

	for _, data := range tabList {

		helper.WriteRowValues(globals.TargetIndexSheet, getMode(data), data.TableName, util.ChangeExtension(data.FileName, ".csv"))
	}

	return nil
}

func getMode(data *helper.MemFileData) (mode string) {
	if data.FileName == "Type.xlsx" {
		mode = "类型表"
	} else {
		mode = "数据表"
	}

	return
}
