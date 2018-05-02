package gosrc

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	err = codegen.NewCodeGen("gosrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(table.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText, globals).
		FormatGoCode().
		WriteBytes(&data).Error()

	return
}
