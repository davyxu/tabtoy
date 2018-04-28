package gosrc

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/davyxu/tabtoy/v3/table"
)

func Generate(globals *model.Globals, fileName string) error {

	return codegen.NewCodeGen("gosrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(table.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText, globals).
		WriteOutputFile(fileName).Error()
}
