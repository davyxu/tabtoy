package model

type Table struct {
	Recs []*Record
}

func (self *Table) Add(r *Record) {
	self.Recs = append(self.Recs, r)
}

// StrStruct: [ {  HP: 1 }, {  HP: 1 } ]
