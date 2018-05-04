package report

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
	"strings"
)

type ErrorObject struct {
	s string

	context []interface{}
}

func getErrorDesc(id string) string {

	errobj := table.CoreConfig.GetKeyValue_ErrorID()
	tobj := reflect.TypeOf(errobj).Elem()
	vobj := reflect.ValueOf(errobj).Elem()

	for i := 0; i < tobj.NumField(); i++ {

		fd := tobj.Field(i)
		if fd.Name == id {
			final := vobj.Field(i).String()
			if final == "" {
				return id
			}

			return final
		}

	}

	return ""
}

func (self *ErrorObject) Error() string {

	var sb strings.Builder

	sb.WriteString(getErrorDesc(self.s))
	sb.WriteString("(")
	sb.WriteString(self.s)
	sb.WriteString(")")
	sb.WriteString(" ")

	for _, c := range self.context {
		sb.WriteString(fmt.Sprintf("%+v", c))
		sb.WriteString(" ")
	}

	return sb.String()
}

func ReportError(s string, context ...interface{}) *ErrorObject {

	panic(&ErrorObject{
		s:       s,
		context: context,
	})
}
