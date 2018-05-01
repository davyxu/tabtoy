package helper

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/table"
	"reflect"
	"strings"
)

type ErrorObject struct {
	s string

	context []interface{}
}

func getErrorDesc(id string) string {

	errobj := table.BuiltinConfig.GetKeyValue_ErrorID()
	tobj := reflect.TypeOf(errobj).Elem()
	vobj := reflect.ValueOf(errobj).Elem()

	for i := 0; i < tobj.NumField(); i++ {

		fd := tobj.Field(i)
		if fd.Name == id {
			return vobj.Field(i).String()
		}

	}

	panic("ErrorID not found: " + id)

}

func (self *ErrorObject) Error() string {

	var sb strings.Builder

	sb.WriteString(getErrorDesc(self.s))
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

func Location(filename string, row, col int) string {

	return fmt.Sprintf("%s(%s)", filename, util.R1C1ToA1(row+1, col+1))
}
