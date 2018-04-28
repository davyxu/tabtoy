package model

type DataTableList struct {
	Data []*DataTable
}

func (self *DataTableList) GetDataTable(name string) *DataTable {

	for _, tab := range self.Data {
		if tab.Name == name {
			return tab
		}
	}

	return nil
}

func (self *DataTableList) AddDataTable(t *DataTable) {
	self.Data = append(self.Data, t)
}
