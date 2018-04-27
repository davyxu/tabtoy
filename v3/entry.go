package v3

import (
	"github.com/davyxu/tabtoy/v3/model"
)

func Parse(globals *model.Globals) error {

	err := loadSymbols(globals, globals.SymbolFile)

	if err != nil {
		return err
	}

	for _, fileName := range globals.InputFileList {

		tab, err := loadTable(fileName)

		if err != nil {
			return err
		}

		resolveHeaderFields(tab, globals.Symbols)

		globals.AddData(tab)
	}

	return nil
}
