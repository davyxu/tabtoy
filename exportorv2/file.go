package exportorv2

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/tealeg/xlsx"
)

type valueRepeatData struct {
	fd    *model.FieldDefine
	value string
}

type File struct {
	TypeSet  *model.BuildInTypeSet // 1个类型描述表
	Name     string                // 电子表格的名称(文件名无后缀)
	FileName string

	valueRepByKey map[valueRepeatData]bool
}

func (self *File) Export(filename string) *model.Table {

	self.Name = makeTableName(filename)
	self.FileName = filename

	log.Infof("%s\n", filepath.Base(filename))

	file, err := xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return nil
	}

	// 解析类型表
	for _, rawSheet := range file.Sheets {

		if isTypeSheet(rawSheet.Name) {
			if self.TypeSet != nil {
				log.Errorln("@Types sheet should keep one!")
				return nil
			}

			typeSheet := newTypeSheet(NewSheet(self, rawSheet))
			if !typeSheet.Parse() {
				return nil
			}

			self.TypeSet = typeSheet.BuildInTypeSet

		}
	}

	// 没有这个类型, 默认填一个
	if self.TypeSet == nil {
		self.TypeSet = model.NewBuildInTypeSet()
	}

	var tab model.Table

	for _, rawSheet := range file.Sheets {

		if !isTypeSheet(rawSheet.Name) {

			log.Infof("            %s", rawSheet.Name)
			dSheet := newDataSheet(NewSheet(self, rawSheet))
			if !dSheet.Export(self, &tab) {
				return nil
			}

		}
	}

	return &tab
}

func (self *File) checkValueRepeat(fd *model.FieldDefine, value string) bool {

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
	return strings.TrimSpace(name) == "@Types"
}

// 制作表名
func makeTableName(filename string) string {
	baseName := path.Base(filename)
	return strings.TrimSuffix(baseName, path.Ext(baseName))
}

func NewFile() *File {

	self := &File{
		valueRepByKey: make(map[valueRepeatData]bool),
	}

	return self
}
