package tests

import (
	"github.com/tealeg/xlsx"
)

type MemFile map[string]*xlsx.File

func (self MemFile) Create(filename string) *xlsx.Sheet {
	sheet := createSheet()

	self[filename] = sheet.File

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
