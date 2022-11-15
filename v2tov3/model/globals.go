package model

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/tealeg/xlsx"
	"path/filepath"
)

type Globals struct {
	TableGetter util.FileGetter

	SourceTypes []ObjectFieldType

	SourceFileList []string

	TargetTypesSheet util.TableSheet

	TargetIndexSheet util.TableSheet

	TargetTables *util.MemFile

	OutputDir string
}

func (self *Globals) AddTableByFile(tableFileName, tableName string, inputFile *xlsx.File) {

	file := util.NewXlsxFile("")

	file.(interface {
		FromXFile(file *xlsx.File)
	}).FromXFile(inputFile)

	tableFileName = filepath.Base(tableFileName)

	self.TargetTables.AddFile(tableFileName, file).TableName = tableName
}

func (self *Globals) AddTable(tableFileName string) util.TableSheet {

	return self.TargetTables.CreateXLSXFile(tableFileName)
}

func (self *Globals) SourceTypeExists(objectTypeName, fieldName string) bool {
	for _, ft := range self.SourceTypes {

		if ft.ObjectType == objectTypeName && ft.FieldName == fieldName {
			return true
		}
	}

	return false
}

func (self *Globals) ObjectTypeByName(objectTypeName string) *ObjectFieldType {
	for _, ft := range self.SourceTypes {

		if ft.ObjectType == objectTypeName {
			return &ft
		}
	}

	return nil
}

func (self *Globals) AddSourceType(oft ObjectFieldType) {

	self.SourceTypes = append(self.SourceTypes, oft)
}

func (self *Globals) PrintTypes() {

	for _, ft := range self.SourceTypes {

		log.Debugf("%+v", ft)
	}
}

func (self *Globals) TypeIsNoneKind(objectTypeName string) bool {
	for _, oft := range self.SourceTypes {
		if oft.ObjectType == objectTypeName && oft.Kind == model.TypeUsage_None {
			return true
		}
	}

	return false
}

func NewGlobals() *Globals {

	return &Globals{
		TargetTables: util.NewMemFile(),
	}
}
