package pbsrc

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

// 将定义用的类型，转换为不同语言对应的复合类型

func init() {
	UsefulFunc["PbType"] = func(tf *model.TypeDefine) string {

		pbType := util.LanguagePrimitive(tf.FieldType, "pb")

		if tf.IsArray() {
			return "repeated " + pbType
		}

		return pbType
	}

	UsefulFunc["PbTag"] = func(fieldIndex int, fieldType *model.TypeDefine) string {

		var sb strings.Builder
		fmt.Fprintf(&sb, "= %d", fieldIndex+1)
		return sb.String()
	}

	UsefulFunc["PbCombineField"] = func(fieldIndex int) string {

		var sb strings.Builder
		fmt.Fprintf(&sb, "= %d", fieldIndex+1)
		return sb.String()
	}
}
