package tests

import (
	"encoding/json"
	"errors"
	"github.com/davyxu/tabtoy/v3"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/jsondata"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

func newGlobal() *model.Globals {
	globals := model.NewGlobals()
	globals.Version = "testver"
	globals.BuiltinSymbolFile = "../table/BuiltinTypes.xlsx"
	globals.IndexFile = "Index.xlsx"
	globals.PackageName = "main"
	globals.CombineStructName = "Config"
	globals.Para = false

	return globals
}

func VerifyType(indexGetter v3.FileGetter, expectJson string) error {

	globals := newGlobal()

	err := v3.Compile(globals, indexGetter)

	if err != nil {
		return err
	}

	appJson := globals.Types.ToJSON()

	globals.Types.Print()

	return compareJson(appJson, []byte(expectJson))
}

func VerifyLauncherJson(indexGetter v3.FileGetter, expectJson string) error {

	globals := newGlobal()

	err := v3.Compile(globals, indexGetter)

	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "tabtoytest_")

	if err != nil {
		return err
	}

	configFileName := filepath.Join(dir, "config.json")

	if err := genFile(globals, configFileName, jsondata.Generate); err != nil {
		return err
	}

	tableFileName := filepath.Join(dir, "table.go")

	if err := genFile(globals, tableFileName, gosrc.Generate); err != nil {
		return err
	}

	appJson, err := compileLauncher(filepath.Join(dir, "launcher.go"), configFileName, tableFileName)
	if err != nil {
		return err
	}

	return compareJson(appJson, []byte(expectJson))
}

func genFile(globals *model.Globals, filename string, genFunc gen.GenFunc) error {

	data, err := genFunc(globals)

	if err != nil {
		return err
	}

	return helper.WriteFile(filename, data)
}

func compareJson(a, b []byte) error {

	var mapA, mapB map[string]interface{}

	err := json.Unmarshal(a, &mapA)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &mapB)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(mapA, mapB) {
		return errors.New("json no match")
	}

	return nil
}
