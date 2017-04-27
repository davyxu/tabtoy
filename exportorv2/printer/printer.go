package printer

type PrinterContext struct {
	outFile string
	p       Printer
	name    string
}

func (self *PrinterContext) Start(g *Globals) bool {

	log.Infof("[%s] %s\n", self.name, self.outFile)

	bf := self.p.Run(g)

	if bf == nil {
		return false
	}

	return bf.WriteFile(self.outFile) == nil
}

type Printer interface {
	Run(g *Globals) *Stream
}

var printerByExt = make(map[string]Printer)

func RegisterPrinter(ext string, p Printer) {

	if _, ok := printerByExt[ext]; ok {
		panic("duplicate printer")
	}

	printerByExt[ext] = p
}
