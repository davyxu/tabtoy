package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func Parse(globals *model.Globals) error {

	err := loadSymbols(globals, globals.BuiltinSymbolFile)

	if err != nil {
		return err
	}
	//
	//if globals.SymbolFile != "" {
	//	err = loadSymbols(globals, globals.SymbolFile)
	//
	//	if err != nil {
	//		return err
	//	}
	//}

	loadPragma(globals, globals.PragmaFile)

	return nil
}
