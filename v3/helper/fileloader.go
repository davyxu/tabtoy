package helper

import (
	"errors"
	"github.com/tealeg/xlsx"
	"sync"
)

type SyncFileLoader struct {
}

func (SyncFileLoader) GetFile(filename string) (*xlsx.File, error) {
	return xlsx.OpenFile(filename)
}

type AsyncFileLoader struct {
	fileByName sync.Map
	inputFile  []string
}

func (self *AsyncFileLoader) AddFile(filename string) {

	self.inputFile = append(self.inputFile, filename)
}

func (self *AsyncFileLoader) Commit() {

	var task sync.WaitGroup
	task.Add(len(self.inputFile))

	for _, inputFileName := range self.inputFile {

		go func(fileName string) {

			file, err := xlsx.OpenFile(fileName)

			if err != nil {
				self.fileByName.Store(fileName, err)
			} else {
				self.fileByName.Store(fileName, file)
			}

			task.Done()

		}(inputFileName)

	}

	task.Wait()

	self.inputFile = self.inputFile[0:0]
}
func (self *AsyncFileLoader) GetFile(filename string) (*xlsx.File, error) {

	if value, ok := self.fileByName.Load(filename); ok {

		switch ret := value.(type) {
		case *xlsx.File:
			return ret, nil
		case error:
			return nil, ret
		default:
			panic("unexpect value")
		}
	} else {
		return nil, errors.New("not found")
	}
}

func NewAsyncFileLoader() *AsyncFileLoader {
	return &AsyncFileLoader{}
}
