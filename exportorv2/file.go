package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/tealeg/xlsx"
)

// 检查单元格值重复结构
type valueRepeatData struct {
	fd    *model.FieldDescriptor
	value string
}

// 1个电子表格文件
type File struct {
	LocalFD  *model.FileDescriptor // 本文件的类型描述表
	GlobalFD *model.FileDescriptor // 全局的类型描述表
	FileName string
	coreFile *xlsx.File

	dataSheets []*DataSheet
	Header     *DataHeader

	valueRepByKey map[valueRepeatData]bool // 检查单元格值重复map
}

func (self *File) ExportLocalType() bool {

	var sheetCount int
	// 解析类型表
	for _, rawSheet := range self.coreFile.Sheets {

		if isTypeSheet(rawSheet.Name) {
			if sheetCount > 0 {
				log.Errorf("%s", i18n.String(i18n.File_TypeSheetKeepSingleton))
				return false
			}

			typeSheet := newTypeSheet(NewSheet(self, rawSheet))

			// 从cell添加类型
			if !typeSheet.Parse(self.LocalFD, self.GlobalFD) {
				return false
			}

			sheetCount++

		}
	}

	// 解析表头
	for _, rawSheet := range self.coreFile.Sheets {

		// 是数据表
		if !isTypeSheet(rawSheet.Name) {
			dSheet := newDataSheet(NewSheet(self, rawSheet))

			if !dSheet.Valid() {
				continue
			}

			log.Infof("            %s", rawSheet.Name)

			dataHeader := newDataHeadSheet()

			// 检查引导头
			if !dataHeader.ParseProtoField(len(self.dataSheets), dSheet.Sheet, self.LocalFD, self.GlobalFD) {
				return false
			}

			if self.Header == nil {
				self.Header = dataHeader
			}

			self.dataSheets = append(self.dataSheets, dSheet)

		}
	}

	return true
}

func (self *File) ExportData() *model.Table {

	self.LocalFD.Name = self.LocalFD.Pragma.TableName

	tab := model.NewTable()
	tab.LocalFD = self.LocalFD

	for _, d := range self.dataSheets {

		log.Infof("            %s", d.Name)

		if !d.Export(self, tab, self.Header) {
			return nil
		}
	}

	return tab
}

func (self *File) checkValueRepeat(fd *model.FieldDescriptor, value string) bool {

	key := valueRepeatData{
		fd:    fd,
		value: value,
	}

	if _, ok := self.valueRepByKey[key]; ok {
		return false
	}

	self.valueRepByKey[key] = true

	return true
}

func isTypeSheet(name string) bool {
	return strings.TrimSpace(name) == model.TypeSheetName
}

func NewFile(filename string) *File {

	self := &File{
		valueRepByKey: make(map[valueRepeatData]bool),
		LocalFD:       model.NewFileDescriptor(),
		FileName:      filename,
	}

	var err error
	self.coreFile, err = xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return nil
	}

	return self
}
