package model

type Globals struct {
	SymbolFile    string
	InputFileList []string
	OutputFile    string

	Symbols SymbolTable

	Datas []*DataTable
}

func (self *Globals) AddData(t *DataTable) {
	self.Datas = append(self.Datas, t)
}
