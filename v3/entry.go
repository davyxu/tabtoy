package v3

import (
	"github.com/davyxu/tabtoy/v3/genfile/json"
	"github.com/davyxu/tabtoy/v3/model"
)

func Run(globals *model.Globals) error {

	globals.OutputFile = "a.json"

	loadSymbols(globals, globals.SymbolFile)

	for _, fileName := range globals.InputFileList {
		loadTable(globals, fileName)
	}

	return json.Generate(globals)
}
