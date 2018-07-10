package model

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/table"
	"github.com/tealeg/xlsx"
	"path/filepath"
)

type Globals struct {
	TableGetter helper.FileGetter

	SourceTypes []ObjectFieldType

	SourceFileList []string

	TargetTypesSheet *xlsx.Sheet

	TargetIndexSheet *xlsx.Sheet

	TargetTables *helper.MemFile

	OutputDir string
}

func (self *Globals) AddTableByFile(tableFileName, tableName string, file *xlsx.File) {

	tableFileName = filepath.Base(tableFileName)

	self.TargetTables.AddFile(tableFileName, file).TableName = tableName
}

func (self *Globals) AddTable(tableFileName, tableName string) *xlsx.File {

	targetFile := xlsx.NewFile()

	tableFileName = filepath.Base(tableFileName)

	self.TargetTables.AddFile(tableFileName, targetFile).TableName = tableName

	return targetFile
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
		if oft.ObjectType == objectTypeName && oft.Kind == table.TableKind_None {
			return true
		}
	}

	return false
}

func NewGlobals() *Globals {

	return &Globals{
		TargetTables: helper.NewMemFile(),
	}
}
