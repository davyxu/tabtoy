package gosrc

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/table"
	"text/template"
)

var UsefulFunc = template.FuncMap{}

func init() {
	UsefulFunc["GoTabTag"] = func(fieldType *table.TypeField) string {
		return fmt.Sprintf("`tab_name:\"%s\"`", fieldType.Name)
	}
}
