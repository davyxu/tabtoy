package model

import (
	"github.com/ahmetb/go-linq"
	"github.com/davyxu/tabtoy/v3/helper"
	"github.com/davyxu/tabtoy/v3/table"
)

type Globals struct {
	Version           string // 工具版本号
	IndexFile         string // 指示文件
	PackageName       string // 文件生成时的包名
	CombineStructName string // 包含最终表所有数据的根结构
	Para              bool   // 并发读取文件

	IndexGetter helper.FileGetter // 索引文件获取器
	TableGetter helper.FileGetter // 其他文件获取器

	IndexList []*table.TablePragma // 输入的索引文件

	Types *TypeTable // 输入的类型及符号

	Datas DataTableList // 输出的字符串格式的数据表

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
