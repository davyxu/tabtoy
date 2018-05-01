package table

func (self *TableField) IsArray() bool {
	return self.ArraySplitter != ""
}
