package helper

import (
	"errors"
	"github.com/davyxu/tabtoy/v3/report"
	"github.com/tealeg/xlsx"
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
}

func (self *FileLoader) AddFile(filename string) {

	self.inputFile = append(self.inputFile, filename)
}

func (self *FileLoader) Commit() {

	var task sync.WaitGroup
	task.Add(len(self.inputFile))

	for _, inputFileName := range self.inputFile {

		go func(fileName string) {

			self.fileByName.Store(fileName, loadFileByExt(fileName))

			task.Done()

		}(inputFileName)

	}

	task.Wait()

	self.inputFile = self.inputFile[0:0]
}

func loadFileByExt(filename string) interface{} {
	switch filepath.Ext(filename) {
	case ".xlsx", ".xls":

		file, err := xlsx.OpenFile(filename)
		if err != nil {
			return err
		}

		file.ToSlice()

		return NewXlsxFile(file)
	case ".csv":

	default:
		report.ReportError("UnknownInputFileExtension", filename)
	}

	return nil
}

func (self *FileLoader) GetFile(filename string) (TableFile, error) {

	if self.syncLoad {

		result := loadFileByExt(filename)
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

func NewFileLoader(syncLoad bool) *FileLoader {
	return &FileLoader{
		syncLoad: syncLoad,
	}
}
