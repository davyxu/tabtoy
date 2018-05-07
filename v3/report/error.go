package report

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
	"strings"
)

type TableError struct {
	ID string

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

func (self *TableError) Error() string {

	var sb strings.Builder

	sb.WriteString("TableError.")
	sb.WriteString(self.ID)
	sb.WriteString(" ")
	sb.WriteString(getErrorDesc(self.ID))
	sb.WriteString(" | ")

	for index, c := range self.context {
		if index > 0 {
			sb.WriteString(" ")
		}

		sb.WriteString(fmt.Sprintf("%+v", c))
	}

	return sb.String()
}

func ReportError(id string, context ...interface{}) *TableError {

	panic(&TableError{
		ID:      id,
		context: context,
	})
}
