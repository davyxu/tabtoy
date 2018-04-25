package json

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/model"
)

func Generate(globals *model.Globals) error {

	return codegen.NewCodeGen("json").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(jsonTeamplate, globals).WriteOutputFile(globals.OutputFile).Error()
}
