package model

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/table"
)

type Globals struct {
	Version           string // 工具版本号
	BuiltinSymbolFile string // 符号文件
	IndexFile         string // 指示文件
	PackageName       string // 文件生成时的包名
	CombineStructName string // 包含最终表所有数据的根结构
	Para              bool   // 并发读取文件

	Types *TypeTable // 类型及符号

	IndexList []*table.TablePragma

	DataTableList // 字符串格式的数据表

	IndexGetter helper.FileGetter
	TableGetter helper.FileGetter
}

func (self *Globals) KeyValueTypeNames() (ret []string) {

	linq.From(self.IndexList).WhereT(func(pragma *table.TablePragma) bool {
		return pragma.TableMode == table.TableMode_KeyValue
	}).SelectT(func(pragma *table.TablePragma) string {

		return pragma.TableType
	}).Distinct().ToSlice(&ret)

	return
}

func NewGlobals() *Globals {
	return &Globals{
		Types: NewSymbolTable(),
	}
}
