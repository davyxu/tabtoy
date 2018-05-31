package model

import (
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/tealeg/xlsx"
)

type Globals struct {
	TableGetter helper.FileGetter

	SourceTypes []ObjectFieldType

	SourceFileList []string

	TargetTypesSheet *xlsx.Sheet

	TargetDatas helper.MemFile
}

func (self *Globals) PrintTypes() {

	for _, ft := range self.SourceTypes {

		log.Debugln(ft.String())
	}
}

func (self *Globals) TypeIsStruct(objectTypeName string) bool {
	for _, oft := range self.SourceTypes {
		if oft.ObjectType == objectTypeName && oft.IsStruct {
			return true
		}
	}

	return false
}
