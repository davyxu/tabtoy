package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func Parse(globals *model.Globals) error {

	loadSymbols(globals, globals.SymbolFile)

	for _, fileName := range globals.InputFileList {
		loadTable(globals, fileName)
	}

	return nil
}
