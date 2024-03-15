package cssrc

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/gen"
	"github.com/davyxu/tabtoy/v4/model"
)

func OutputFile(globals *model.Globals, outFile string) (err error) {

	cg := codegen.NewCodeGen("cssrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc)

	err = cg.ParseTemplate(templateText, globals).Error()
	if err != nil {
		return
	}

	return util.WriteFile(outFile, cg.Data())
}
