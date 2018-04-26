package model

type Globals struct {
	Version       string
	SymbolFile    string
	InputFileList []string
	PackageName   string

	Symbols SymbolTable

	Datas []*DataTable
}

func (self *Globals) AddData(t *DataTable) {
	self.Datas = append(self.Datas, t)
}
