package model

import (
	"github.com/davyxu/tabtoy/util"
)

type Compiler struct {
	DataFileGetter util.FileGetter // 数据源
	ParaLoading    bool
	CacheDir       string

	FileList []string // 输入的文件列表

	Types *TypeTable // 输入的类型及符号

	Datas DataTableList
}

func NewCompiler() *Compiler {
	return &Compiler{
		Types: NewTypeTable(),
	}
}
