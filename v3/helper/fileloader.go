package helper

import (
	"errors"
	"github.com/davyxu/tabtoy/v3/report"
	"path/filepath"
	"sync"
)

type FileGetter interface {
	GetFile(filename string) (TableFile, error)
}

type FileLoader struct {
	fileByName sync.Map
	inputFile  []string

	syncLoad bool
	cacheDir string

	UseGBKCSV bool
}

func (self *FileLoader) AddFile(filename string) {

	self.inputFile = append(self.inputFile, filename)
}

func (self *FileLoader) Commit() {

	var task sync.WaitGroup
	task.Add(len(self.inputFile))

	for _, inputFileName := range self.inputFile {

		go func(fileName string) {

			self.fileByName.Store(fileName, loadFileByExt(fileName, self.UseGBKCSV, self.cacheDir))

			task.Done()

		}(inputFileName)

	}

	task.Wait()

	self.inputFile = self.inputFile[0:0]
}

func loadFileByExt(filename string, useGBKCSV bool, cacheDir string) interface{} {

	var tabFile TableFile
	switch filepath.Ext(filename) {
	case ".xlsx", ".xls", ".xlsm":

		tabFile = NewXlsxFile(cacheDir)

		err := tabFile.Load(filename)

		if err != nil {
			return err
		}

	case ".csv":
		tabFile = NewCSVFile()

		err := tabFile.Load(filename)

		if err != nil {
			return err
		}

		// 输入gbk, 内部utf8
		if useGBKCSV {
			tabFile.(*CSVFile).Transform(ConvGBKToUTF8)
		}

	default:
		report.ReportError("UnknownInputFileExtension", filename)
	}

	return tabFile
}

func (self *FileLoader) GetFile(filename string) (TableFile, error) {

	if self.syncLoad {

		result := loadFileByExt(filename, self.UseGBKCSV, self.cacheDir)
		if err, ok := result.(error); ok {
			return nil, err
		}

		return result.(TableFile), nil

	} else {
		if result, ok := self.fileByName.Load(filename); ok {

			if err, ok := result.(error); ok {
				return nil, err
			}

			return result.(TableFile), nil

		} else {
			return nil, errors.New("not found")
		}
	}

}

func NewFileLoader(syncLoad bool, cacheDir string) *FileLoader {
	return &FileLoader{
		syncLoad: syncLoad,
		cacheDir: cacheDir,
	}
}
