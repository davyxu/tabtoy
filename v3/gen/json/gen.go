package json

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/model"
)

func Generate(globals *model.Globals, fileName string) error {

	return codegen.NewCodeGen("json").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText, globals).
		WriteOutputFile(fileName).Error()
}
