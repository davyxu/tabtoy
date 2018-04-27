package gosrc

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/table"
	"strings"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func init() {
	UsefulFunc["GoTabTag"] = func(fieldType *table.TableField) string {

		var sb strings.Builder

		sb.WriteString("`")

		sb.WriteString(fmt.Sprintf("tb_name:\"%s\"", fieldType.Name))
		//
		//if fieldType.Splitter != "" {
		//	sb.WriteString(" ")
		//	sb.WriteString(fmt.Sprintf("tb_splitter:\"%s\"", fieldType.Splitter))
		//}

		sb.WriteString("`")

		return sb.String()
	}
}
