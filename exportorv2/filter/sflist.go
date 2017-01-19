package filter

import (
	"sort"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

type structFieldResult struct {
	key   *model.FieldDescriptor
	value string
}

type structFieldList struct {
	data []*structFieldResult
}

func (self *structFieldList) Add(fd *model.FieldDescriptor, value string) {

	self.data = append(self.data, &structFieldResult{
		key:   fd,
		value: value,
	})

}

func (self *structFieldList) Exists(fd *model.FieldDescriptor) bool {
	for _, libfd := range self.data {
		if libfd.key == fd {
			return true
		}
	}

	return false
}

func (self *structFieldList) Len() int {
	return len(self.data)
}

func (self *structFieldList) Get(index int) *structFieldResult {
	return self.data[index]
}

func (self *structFieldList) Swap(i, j int) {
	self.data[i], self.data[j] = self.data[j], self.data[i]
}

func (self *structFieldList) Less(i, j int) bool {

	a := self.data[i]
	b := self.data[j]

	return a.key.Order < b.key.Order
}

func (self *structFieldList) Sort() {
	sort.Sort(self)
}

func newStructFieldList() *structFieldList {
	return &structFieldList{}
}
