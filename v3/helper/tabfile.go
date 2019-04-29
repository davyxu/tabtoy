package helper

type TableFile interface {

	// 获取所有表单
	Sheets() []TableSheet
}

type TableSheet interface {

	// 表单名称
	Name() string

	// 从表单指定单元格获取值
	GetValue(row, col int, isFloat bool) string

	// 最大列
	MaxColumn() int
}

// 检查表单的某行是否全空
func IsRowEmpty(sheet TableSheet, row int) bool {

	for col := 0; col < sheet.MaxColumn(); col++ {

		data := sheet.GetValue(row, col, false)

		if data != "" {
			return false
		}
	}

	return true
}
