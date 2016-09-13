package printer

import (
	"path/filepath"
)

type PrinterContext struct {
	outFile string
	p       Printer
	ext     string
}

func (self *PrinterContext) Start(g *Globals) bool {
	filebase := filepath.Base(self.outFile)

	log.Infof("%s\n", filebase)

	bf := self.p.Run(g)

	if bf == nil {
		return false
	}

	return bf.Write(self.outFile)
}

type Printer interface {
	Run(g *Globals) *BinaryFile
}

var printerByExt = make(map[string]Printer)

func RegisterPrinter(ext string, p Printer) {

	if _, ok := printerByExt[ext]; ok {
		panic("duplicate printer")
	}

	printerByExt[ext] = p
}
