package printer

type modlistPrinter struct {
}

func (self *modlistPrinter) Run(g *Globals) *Stream {

	bf := NewStream()

	for _, filename := range g.ModList {
		bf.Printf("%s\n", filename)
	}

	return bf
}

func init() {

	RegisterPrinter("modlist", &modlistPrinter{})

}
