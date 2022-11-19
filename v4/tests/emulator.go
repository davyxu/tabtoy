package tests

import (
	"encoding/json"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/compiler"
	"github.com/davyxu/tabtoy/v4/gen/gosrc"
	"github.com/davyxu/tabtoy/v4/gen/jsondata"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func init() {
	report.Init()
}

type TableEmulator struct {
	G *model.Globals
	T *testing.T

	memFile *util.MemFile

	compiled bool
}

func NewTableEmulator(t *testing.T) *TableEmulator {

	memFile := util.NewMemFile()
	g := model.NewGlobals()
	g.DataFileGetter = memFile

	return &TableEmulator{
		T:       t,
		G:       g,
		memFile: memFile,
	}
}

func (self *TableEmulator) SetSheetName(sheet util.TableSheet, name string) {
	if csvSheet, ok := sheet.(*util.CSVSheet); ok {
		csvSheet.SetName(name)
	}
}

func (self *TableEmulator) createSheet(mode, fileName string) util.TableSheet {

	sheet := self.memFile.CreateCSVFile(fileName)
	self.G.AddFile(mode, fileName)
	return sheet
}

func (self *TableEmulator) CreateDataSheet(fileName string) util.TableSheet {
	return self.createSheet("", fileName)
}

func (self *TableEmulator) CreateKVSheet(fileName string) util.TableSheet {
	sheet := self.createSheet("KV", fileName)
	sheet.WriteRow("Key", "Type", "Value", "Comment", "Meta")
	return sheet
}

func (self *TableEmulator) run(errThenFail bool) error {
	if self.compiled {
		return nil
	}

	runtime.GOMAXPROCS(1)

	var err error

	if errThenFail {
		defer func() {
			if err != nil {
				self.T.Fatal(err)
			}

		}()
	}

	err = compiler.Compile(self.G)

	self.compiled = true

	return err
}

func (self *TableEmulator) MustGotError(expectError string) {

	err := self.run(false)

	if err == nil {
		self.T.Fatalf("Expect error '%s' but no error", expectError)
	}

	nowErr := err.Error()
	if nowErr != expectError {
		self.T.Fatalf("Expect '%s' \ngot '%s'", expectError, nowErr)
	}
}

func (self *TableEmulator) VerifyType(expectJson string) {

	self.run(true)

	appJson := self.G.Types.ToJSON()

	result, _ := compareArrayJson(appJson, []byte(expectJson))

	if !result {
		self.T.Fatalf("Expect '%s' got '%s'", expectJson, appJson)
	}
}

func (self *TableEmulator) VerifyData(expectJson string) {

	self.run(true)

	appJson, err := jsondata.OutputData(self.G)

	if err != nil {
		return
	}

	var result bool
	result, err = compareKVJson(appJson, []byte(expectJson))

	if !result {
		self.T.Fatalf("Expect '%s' got '%s'", expectJson, appJson)
	}
}

func (self *TableEmulator) VerifyGoTypeAndJson(expectJson string) {

	self.run(true)

	dir, err := os.MkdirTemp("", "tabtoytest_")

	if err != nil {
		self.T.Fatal(err)
		return
	}

	configFileName := filepath.Join(dir, "config.json")

	if err = jsondata.OutputFile(self.G, configFileName); err != nil {
		return
	}

	tableFileName := filepath.Join(dir, "table.go")

	if err = gosrc.OutputFile(self.G, tableFileName); err != nil {
		return
	}

	var appJson []byte
	appJson, err = compileLauncher(filepath.Join(dir, "launcher.go"), configFileName, tableFileName)
	if err != nil {
		self.T.Fatal(string(appJson))
		return
	}

	var result bool
	result, err = compareKVJson(appJson, []byte(expectJson))
	if err != nil {
		return
	}

	if !result {
		self.T.Fatalf("Expect '%s' got '%s'", expectJson, appJson)
	}
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
