package model

type Globals struct {
	Version           string
	SymbolFile        string
	InputFileList     []string
	PackageName       string
	CombineStructName string // 包含最终表所有数据的根结构

	Symbols SymbolTable

	Datas []*DataTable
}

func (self *Globals) AddData(t *DataTable) {
	self.Datas = append(self.Datas, t)
}
