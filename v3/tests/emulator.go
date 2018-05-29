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
	"testing"
)

type TableEmulator struct {
	G *model.Globals
	T *testing.T

	helper.MemFile
}

func NewTableEmulator(t *testing.T) *TableEmulator {

	globals := model.NewGlobals()
	globals.Version = "testver"
	globals.BuiltinSymbolFile = "../table/BuiltinTypes.xlsx"
	globals.IndexFile = "Index.xlsx"
	globals.PackageName = "main"
	globals.CombineStructName = "Config"
	globals.Para = false

	memfile := helper.NewMemFile()

	globals.TableGetter = memfile
	globals.IndexGetter = memfile

	return &TableEmulator{G: globals, T: t, MemFile: memfile}
}

func (self *TableEmulator) VerifyError(expectError string) {

	err := v3.Compile(self.G)

	if err == nil || err.Error() != expectError {
		self.T.Logf("Expect '%s' got '%s'", expectError, err.Error())
		self.T.FailNow()
	}
}

func (self *TableEmulator) VerifyType(expectJson string) error {

	err := v3.Compile(self.G)

	if err != nil {
		return err
	}

	appJson := self.G.Types.ToJSON()

	println(string(appJson))

	if expectJson == "" {
		return nil
	}

	return compareJson(appJson, []byte(expectJson))
}

func (self *TableEmulator) VerifyLauncherJson(expectJson string) error {

	err := v3.Compile(self.G)

	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir("", "tabtoytest_")

	if err != nil {
		return err
	}

	configFileName := filepath.Join(dir, "config.json")

	if err := genFile(self.G, configFileName, jsondata.Generate); err != nil {
		return err
	}

	tableFileName := filepath.Join(dir, "table.go")

	if err := genFile(self.G, tableFileName, gosrc.Generate); err != nil {
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
