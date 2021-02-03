package luasrc

import (
	"fmt"
	"github.com/davyxu/protoplus/codegen"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	err = codegen.NewCodeGen("luasrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText_luasrc, globals).
		WriteBytes(&data).Error()

	return
}

func Output(globals *model.Globals, param string) (err error) {

	type LocalContext struct {
		Tab *model.DataTable
		G   *model.Globals
	}

	var typeData []byte
	err = codegen.NewCodeGen("luatype").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc).
		ParseTemplate(templateText_luatype, globals).
		WriteBytes(&typeData).Error()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/_%sType.lua", param, globals.CombineStructName), typeData, 0666)

	if err != nil {
		return err
	}

	for _, tab := range globals.Datas.AllTables() {

		var ctx LocalContext
		ctx.Tab = tab
		ctx.G = globals

		var data []byte
		err = codegen.NewCodeGen("luadir").
			RegisterTemplateFunc(codegen.UsefulFunc).
			RegisterTemplateFunc(gen.UsefulFunc).
			RegisterTemplateFunc(UsefulFunc).
			ParseTemplate(templateText_luadir, ctx).
			WriteBytes(&data).Error()

		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.lua", param, tab.HeaderType), data, 0666)

		if err != nil {
			return err
		}
	}

	return nil
}
