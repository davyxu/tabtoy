package model

type DataTableList struct {
	Data []*DataTable
}

func (self *DataTableList) GetDataTable(headerType string) *DataTable {

	for _, tab := range self.Data {
		if tab.HeaderType == headerType {
			return tab
		}
	}

	return nil
}

func (self *DataTableList) AddDataTable(t *DataTable) {
	self.Data = append(self.Data, t)
}
