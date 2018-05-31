package model

type DataTableList struct {
	data []*DataTable
}

func (self *DataTableList) GetDataTable(headerType string) *DataTable {

	for _, tab := range self.data {
		if tab.HeaderType == headerType {
			return tab
		}
	}

	return nil
}

func (self *DataTableList) AddDataTable(t *DataTable) {
	self.data = append(self.data, t)
}
func (self *DataTableList) AllTables() []*DataTable {
	return self.data
}
