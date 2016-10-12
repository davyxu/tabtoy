// Generated by github.com/davyxu/tabtoy
// Version: 2.2.2
// DO NOT EDIT!!
package table

import (
	"gamedef"
	"fmt"
)

var SampleByID = make(map[int64]*gamedef.ItemDefine)

var SampleByName = make(map[string]*gamedef.ItemDefine)

var ExpByLevel = make(map[int32]*gamedef.ItemDefine)

func MakeConfigIndex(v *gamedef.Config) {

	// Sample
	for _, def := range v.Sample {

		if _, ok := SampleByID[def.ID]; ok {
			panic(fmt.Sprintf("duplicate index in SampleByID: %v", def.ID))
		}

		if _, ok := SampleByName[def.Name]; ok {
			panic(fmt.Sprintf("duplicate index in SampleByName: %v", def.Name))
		}

		SampleByID[def.ID] = def
		SampleByName[def.Name] = def
	}

	// Exp
	for _, def := range v.Exp {

		if _, ok := ExpByLevel[def.Level]; ok {
			panic(fmt.Sprintf("duplicate index in ExpByLevel: %v", def.Level))
		}

		ExpByLevel[def.Level] = def
	}

}
