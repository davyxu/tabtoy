package tests

import (
	"github.com/davyxu/tabtoy/v3"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/json"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
	"path/filepath"
)

func Run(indexGetter v3.FileGetter) error {

	globals := model.NewGlobals()
	globals.Version = "testver"
	globals.BuiltinSymbolFile = "../table/BuiltinTypes.xlsx"
	globals.IndexFile = "Index.xlsx"
	globals.PackageName = "main"
	globals.CombineStructName = "Config"
	globals.Para = false

	err := v3.Compile(globals, indexGetter)

	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "tabtoytest_")

	if err != nil {
		return err
	}

	configFileName := filepath.Join(dir, "config.json")

	if err := genFile(globals, configFileName, json.Generate); err != nil {
		return err
	}

	tableFileName := filepath.Join(dir, "table.go")

	if err := genFile(globals, tableFileName, gosrc.Generate); err != nil {
		return err
	}

	return compileLauncher(filepath.Join(dir, "launcher.go"), configFileName, tableFileName)
}

func genFile(globals *model.Globals, filename string, genFunc gen.GenFunc) error {

	data, err := genFunc(globals)

	if err != nil {
		return err
	}

	return helper.WriteFile(filename, data)
}
