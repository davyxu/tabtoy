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

	mergeList []*File
}

func (self *File) GlobalFileDesc() *model.FileDescriptor {
	return self.GlobalFD

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

	// 表头用到的表单
	var headerSheet *DataSheet

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
				headerSheet = dSheet
			} else {
				if !self.Header.Equal(dataHeader) {
					log.Errorf("%s %s!=%s", i18n.String(i18n.DataHeader_NotMatch), headerSheet.Name, dSheet.Name)
					return false
				}
			}

			self.dataSheets = append(self.dataSheets, dSheet)

		}
	}

	// File描述符的名字必须放在类型里, 因为这里始终会被调用, 但是如果数据表缺失, 是不会更新Name的
	self.LocalFD.Name = self.LocalFD.Pragma.GetString("TableName")

	return true
}

func (self *File) IsVertical() bool {
	return self.LocalFD.Pragma.GetBool("Vertical")
}

func (self *File) ExportData(dataModel *model.DataModel, parentHeader *DataHeader) bool {

	for _, d := range self.dataSheets {

		log.Infof("            %s", d.Name)

		if !d.Export(self, dataModel, self.Header, parentHeader) {
			return false
		}
	}

	return true

}

func (self *File) CheckValueRepeat(fd *model.FieldDescriptor, value string) bool {

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
		log.Errorf("%s, %v", i18n.String(i18n.System_OpenReadXlsxFailed), err.Error())

		return nil
	}

	return self
}
