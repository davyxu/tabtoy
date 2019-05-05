package helper

import (
	"errors"
	"github.com/tealeg/xlsx"
)

type MemFileData struct {
	File      TableFile
	TableName string
	FileName  string
}

type MemFile struct {
	dataByFileName map[string]*MemFileData
}

func (self *MemFile) VisitAllTable(callback func(data *MemFileData) bool) {

	for _, v := range self.dataByFileName {
		if !callback(v) {
			return
		}
	}
}

func (self *MemFile) AddFile(filename string, file TableFile) (ret *MemFileData) {

	ret = &MemFileData{
		File:     file,
		FileName: filename,
	}

	self.dataByFileName[filename] = ret

	return
}

func (self *MemFile) CreateXLSXFile(filename string) TableSheet {

	xfile := xlsx.NewFile()
	xfile.AddSheet("Default")

	file := NewXlsxFile()

	file.(interface {
		FromXFile(file *xlsx.File)
	}).FromXFile(xfile)

	self.AddFile(filename, file)

	return file.Sheets()[0]
}

func (self *MemFile) CreateCSVFile(filename string) TableSheet {

	file := NewCSVFile()

	self.AddFile(filename, file)

	return file.Sheets()[0]
}

func (self *MemFile) GetFile(filename string) (TableFile, error) {

	if f, ok := self.dataByFileName[filename]; ok {

		return f.File, nil
	}

	return nil, errors.New("file not found: " + filename)
}

func NewMemFile() *MemFile {
	return &MemFile{
		dataByFileName: make(map[string]*MemFileData),
	}
}
