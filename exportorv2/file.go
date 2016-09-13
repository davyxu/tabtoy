package exportorv2

import (
	"path/filepath"
	"strings"

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
	*model.FileDescriptor // 1个类型描述表
	FileName              string

	valueRepByKey map[valueRepeatData]bool // 检查单元格值重复map
}

func (self *File) readTypeSheet(file *xlsx.File) bool {

	var sheetCount int
	// 解析类型表
	for _, rawSheet := range file.Sheets {

		if isTypeSheet(rawSheet.Name) {
			if sheetCount > 0 {
				log.Errorf("%s sheet should keep one!", model.TypeSheetName)
				return false
			}

			typeSheet := newTypeSheet(NewSheet(self, rawSheet))

			// 从cell添加类型
			if !typeSheet.Parse(self.FileDescriptor) {
				return false
			}

			sheetCount++

		}
	}

	return true
}

func (self *File) Export(filename string) *model.Table {

	self.FileName = filename

	log.Infof("%s\n", filepath.Base(filename))

	file, err := xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return nil
	}

	// 读取类型表
	if !self.readTypeSheet(file) {
		return nil
	}

	self.Name = self.Pragma.TableName

	tab := model.NewTable(self.FileDescriptor)

	// 遍历数据表
	for _, rawSheet := range file.Sheets {

		if !isTypeSheet(rawSheet.Name) {

			dSheet := newDataSheet(NewSheet(self, rawSheet))

			if !dSheet.Valid() {
				continue
			}

			log.Infof("            %s", rawSheet.Name)

			dataHeader := newDataHeadSheet()

			// 检查引导头
			if !dataHeader.ParseProtoField(dSheet.Sheet, self.FileDescriptor) {
				return nil
			}

			// TODO 只使用第一个sheet中的protoheader定义
			// TODO 其他Sheet可以在顶部定义一个标记@RefProtoHeader, 引用前面的protoheader
			if !dSheet.Export(self, tab, dataHeader) {
				return nil
			}

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

func NewFile() *File {

	self := &File{
		valueRepByKey:  make(map[valueRepeatData]bool),
		FileDescriptor: model.NewFileDescriptor(),
	}

	return self
}
