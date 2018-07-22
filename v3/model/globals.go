package model

import (
	"github.com/davyxu/tabtoy/v3/helper"
)

type Globals struct {
	Version           string // 工具版本号
	IndexFile         string // 指示文件
	PackageName       string // 文件生成时的包名
	CombineStructName string // 包含最终表所有数据的根结构
	//Para              bool   // 并发读取文件

	IndexGetter helper.FileGetter // 索引文件获取器
	TableGetter helper.FileGetter // 其他文件获取器

	IndexList []*IndexDefine // 输入的索引文件

	Types *TypeTable // 输入的类型及符号

	Datas DataTableList // 输出的字符串格式的数据表
}

func NewGlobals() *Globals {
	return &Globals{
		Types: NewSymbolTable(),
	}
}
