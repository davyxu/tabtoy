package model

import (
	"github.com/davyxu/tabtoy/util"
)

type Globals struct {
	Version           string // 工具版本号
	PackageName       string // 文件生成时的包名
	CombineStructName string // 包含最终表所有数据的根结构

	IndexGetter util.FileGetter // 索引文件获取器
	TableGetter util.FileGetter // 其他文件获取器

	GenBinary bool

	ParaLoading bool

	CacheDir string
}

func NewGlobals() *Globals {
	return &Globals{}
}
