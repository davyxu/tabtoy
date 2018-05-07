package model

import (
	"fmt"
	"github.com/davyxu/tabtoy/util"
)

type Cell struct {
	Value string
	Row   int // base 0
	Col   int // base 0
	File  string
	Sheet string
}

func (self *Cell) String() string {
	return fmt.Sprintf("'%s' @%s|%s(%s)", self.Value, self.File, self.Sheet, util.R1C1ToA1(self.Row+1, self.Col+1))
}
