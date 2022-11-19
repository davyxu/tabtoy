package gosrc

import (
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/gen"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
)

func OutputFile(globals *model.Globals, outFile string) (err error) {

	cg := codegen.NewCodeGen("gosrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc)

	err = cg.ParseTemplate(templateText, globals).Error()
	if err != nil {
		return
	}

	err = cg.FormatGoCode().Error()
	if err != nil {
		report.Log.Infoln(cg.Code())
		return
	}

	return util.WriteFile(outFile, cg.Data())
}
