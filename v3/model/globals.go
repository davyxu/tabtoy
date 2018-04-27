package model

type Globals struct {
	Version           string   // 工具版本号
	BuiltinSymbolFile string   // 符号文件
	SymbolFile        string   // 符号文件
	PragmaFile        string   // 指示文件
	InputFileList     []string // 输入的电子表格文件列表
	PackageName       string   // 文件生成时的包名
	CombineStructName string   // 包含最终表所有数据的根结构

	Symbols SymbolTable // 类型及符号

	Datas []*DataTable // 字符串格式的数据
}

func (self *Globals) AddData(t *DataTable) {
	self.Datas = append(self.Datas, t)
}
