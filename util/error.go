package util

import (
	"fmt"
	"strings"
)

type TableError struct {
	ID string

	context []interface{}
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

type ErrorLanguage struct {
	CHS string
}

var (
	errorByID = map[string]*ErrorLanguage{}
)

func getErrorDesc(id string) string {

	if lan, ok := errorByID[id]; ok {
		return lan.CHS
	}

	return ""
}
func RegisterError(key string, lan *ErrorLanguage) {
	if _, ok := errorByID[key]; ok {
		panic("duplicate error define: " + key)
	}
	errorByID[key] = lan
}
