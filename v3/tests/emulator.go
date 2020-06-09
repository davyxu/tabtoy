package tests

import (
	"encoding/json"
	"github.com/davyxu/tabtoy/v3/compiler"
	"github.com/davyxu/tabtoy/v3/gen"
	"github.com/davyxu/tabtoy/v3/gen/gosrc"
	"github.com/davyxu/tabtoy/v3/gen/jsontext"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/model"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

type TableEmulator struct {
	G *model.Globals
	T *testing.T

	*helper.MemFile
}

func NewTableEmulator(t *testing.T) *TableEmulator {

	globals := model.NewGlobals()
	globals.Version = "testver"
	globals.IndexFile = "Index"
	globals.PackageName = "main"
	globals.CombineStructName = "Table"

	memfile := helper.NewMemFile()

	globals.TableGetter = memfile
	globals.IndexGetter = memfile

	return &TableEmulator{
		G:       globals,
		T:       t,
		MemFile: memfile}
}

func (self *TableEmulator) MustGotError(expectError string) {

	err := compiler.Compile(self.G)

	if err == nil {
		self.T.Logf("Expect error '%s'", expectError)
		self.T.FailNow()
	}

	if err.Error() != expectError {
		self.T.Logf("Expect '%s' got '%s'", expectError, err.Error())
		self.T.FailNow()
	}
}

func (self *TableEmulator) VerifyType(expectJson string) {

	runtime.GOMAXPROCS(1)

	var err error

	defer func() {
		if err != nil {
			self.T.Error(err)
			self.T.FailNow()
		}

	}()

	err = compiler.Compile(self.G)

	if err != nil {
		return
	}

	appJson := self.G.Types.ToJSON()

	if expectJson == "" {
		return
	}

	var result bool
	result, err = compareArrayJson(appJson, []byte(expectJson))
	if err != nil {
		self.T.Logf("InvalidJson '%s'", string(appJson))
		return
	}

	if !result {
		self.T.Logf("Expect '%s' got '%s'", expectJson, appJson)
		self.T.FailNow()
	}
}

func (self *TableEmulator) VerifyData(expectJson string) {

	runtime.GOMAXPROCS(1)

	var err error

	defer func() {
		if err != nil {
			self.T.Error(err)
			self.T.FailNow()
		}

	}()

	err = compiler.Compile(self.G)

	if err != nil {
		return
	}

	var appJson []byte
	appJson, err = jsondata.Generate(self.G)

	if err != nil {
		return
	}

	var result bool
	result, err = compareKVJson(appJson, []byte(expectJson))
	if err != nil {
		self.T.Logf("InvalidJson '%s'", string(appJson))
		return
	}

	if !result {
		self.T.Logf("Expect '%s' got '%s'", expectJson, appJson)
		self.T.FailNow()
	}
}

func (self *TableEmulator) VerifyGoTypeAndJson(expectJson string) {

	var err error

	defer func() {
		if err != nil {
			self.T.Error(err)
			self.T.FailNow()
		}

	}()

	err = compiler.Compile(self.G)

	if err != nil {
		return
	}

	var dir string
	dir, err = ioutil.TempDir("", "tabtoytest_")

	if err != nil {
		return
	}

	configFileName := filepath.Join(dir, "config.json")

	if err = genFile(self.G, configFileName, jsondata.Generate); err != nil {
		return
	}

	tableFileName := filepath.Join(dir, "table.go")

	if err = genFile(self.G, tableFileName, gosrc.Generate); err != nil {
		return
	}

	var appJson []byte
	appJson, err = compileLauncher(filepath.Join(dir, "launcher.go"), configFileName, tableFileName)
	if err != nil {
		self.T.Log(string(appJson))
		self.T.FailNow()
		return
	}

	var result bool
	result, err = compareKVJson(appJson, []byte(expectJson))
	if err != nil {
		return
	}

	if !result {
		self.T.Logf("Expect '%s' got '%s'", expectJson, appJson)
		self.T.FailNow()
	}
}

func genFile(globals *model.Globals, filename string, genFunc gen.GenFunc) error {

	data, err := genFunc(globals)

	if err != nil {
		return err
	}

	return helper.WriteFile(filename, data)
}

func compareKVJson(a, b []byte) (bool, error) {

	var mapA, mapB map[string]interface{}

	err := json.Unmarshal(a, &mapA)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(b, &mapB)
	if err != nil {
		return false, err
	}

	if !reflect.DeepEqual(mapA, mapB) {
		return false, nil
	}

	return true, nil
}

func compareArrayJson(a, b []byte) (bool, error) {

	var mapA, mapB []interface{}

	err := json.Unmarshal(a, &mapA)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(b, &mapB)
	if err != nil {
		return false, err
	}

	if !reflect.DeepEqual(mapA, mapB) {
		return false, nil
	}

	return true, nil
}
