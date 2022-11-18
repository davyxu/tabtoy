package model

type TableManager struct {
	data []*DataTable
}

func (self *TableManager) GetDataTable(headerType string) *DataTable {

	for _, tab := range self.data {
		if tab.HeaderType == headerType {
			return tab
		}
	}

	return nil
}

func (self *TableManager) AddDataTable(t *DataTable) {
	self.data = append(self.data, t)
}
func (self *TableManager) AllTables() []*DataTable {
	return self.data
}

func (self *TableManager) Count() int {
	return len(self.data)
}

func NewTableManager() *TableManager {
	return &TableManager{}
}
