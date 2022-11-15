package main

import (
	"flag"
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2tov3"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"os"
)

var (
	paramUpgradeOut = flag.String("up_out", "", "upgrade v2 table to v3 format output dir")
)

func V2ToV3Entry() {

	globals := model.NewGlobals()

	globals.TableGetter = util.NewFileLoader(true, "")

	globals.SourceFileList = flag.Args()
	globals.OutputDir = *paramUpgradeOut

	if err := v2tov3.Upgrade(globals); err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

}
