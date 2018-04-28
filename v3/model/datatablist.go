package model

type DataTableList struct {
	Datas []*DataTable
}

func (self *DataTableList) GetDataTable(name string) *DataTable {

	for _, tab := range self.Datas {
		if tab.Name == name {
			return tab
		}
	}

	return nil
}

func (self *DataTableList) AddDataTable(t *DataTable) {
	self.Datas = append(self.Datas, t)
}
