package gen

import (
	"github.com/davyxu/tabtoy/v3/model"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

type TableIndices struct {
	Table     *model.DataTable
	FieldInfo *model.TypeDefine
}

// 将定义用的类型，转换为不同语言对应的复合类型
func LanguageType(tf *model.TypeDefine, lanType string) string {

	convertedType := model.LanguagePrimitive(tf.FieldType, lanType)

	if tf.IsArray() {
		switch lanType {
		case "cs":
			return convertedType + "[]"
		case "go":
			return "[]" + convertedType
		default:
			panic("unknown lan type: " + lanType)
		}
	}

	return convertedType
}

func init() {
	UsefulFunc["LanguageType"] = LanguageType

	UsefulFunc["GetIndices"] = func(globals *model.Globals) (ret []TableIndices) {

		for _, tab := range globals.Datas.AllTables() {

			// 遍历输入数据的每一列
			for _, header := range tab.Headers {

				// 输入的列头
				if header.TypeInfo == nil {
					continue
				}

				if header.TypeInfo.MakeIndex {

					ret = append(ret, TableIndices{
						Table:     tab,
						FieldInfo: header.TypeInfo,
					})
				}
			}
		}

		return

	}
}
