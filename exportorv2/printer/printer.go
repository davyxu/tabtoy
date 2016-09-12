package printer

import (
	"path"
)

type PrinterContext struct {
	outDir string
	p      Printer
	ext    string
}

func (self *PrinterContext) Start(g *Globals) bool {
	filebase := g.CombineStructName + self.ext
	outputFile := path.Join(self.outDir, filebase)

	log.Infof("%s\n", filebase)

	bf := self.p.Run(g)

	if bf == nil {
		return false
	}

	return bf.Write(outputFile)
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
