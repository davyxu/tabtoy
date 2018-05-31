package helper

import (
	"github.com/tealeg/xlsx"
)

type MemFile map[string]*xlsx.File

func (self MemFile) AddFile(filename string, file *xlsx.File) {
	self[filename] = file
}

func (self MemFile) Create(filename string) *xlsx.Sheet {

	f := xlsx.NewFile()
	sheet, _ := f.AddSheet("Default")

	self.AddFile(filename, sheet.File)

	return sheet
}

func (self MemFile) GetFile(filename string) (*xlsx.File, error) {

	if f, ok := self[filename]; ok {
		return f, nil
	}

	return xlsx.OpenFile(filename)
}

func NewMemFile() MemFile {
	return make(MemFile)
}
