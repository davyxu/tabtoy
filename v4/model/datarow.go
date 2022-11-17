package model

type DataRow struct {
	row   int
	cells []*Cell
	tab   *DataTable
}

func (self *DataRow) Cells() []*Cell {
	return self.cells
}

func (self *DataRow) Cell(col int) *Cell {
	return self.cells[col]
}

func (self *DataRow) AddCell() (ret *Cell) {

	ret = &Cell{
		Col:   len(self.cells),
		Row:   self.row,
		Table: self.tab,
	}

	self.cells = append(self.cells, ret)
	return
}

func (self *DataRow) IsEmpty() bool {
	return len(self.cells) == 0
}

func newDataRow(row int, tab *DataTable) *DataRow {
	return &DataRow{
		row: row,
		tab: tab,
	}
}
