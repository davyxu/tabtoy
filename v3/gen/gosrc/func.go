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

		var kv []string

		if fieldType.Name != "" {
			kv = append(kv, fmt.Sprintf("tb_name:\"%s\"", fieldType.Name))
		}

		if len(kv) > 0 {
			sb.WriteString("`")

			for _, s := range kv {
				sb.WriteString(s)
			}

			sb.WriteString("`")
		}

		return sb.String()
	}
}
