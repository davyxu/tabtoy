package main

import (
	"flag"
	"github.com/davyxu/tabtoy/v2tov3"
	"github.com/davyxu/tabtoy/v2tov3/model"
	"github.com/davyxu/tabtoy/v3/helper"
	"os"
	"path/filepath"
)

var (
	paramUpgradeOut = flag.String("upout", "", "upgrade v2 table to v3 format output dir")
)

func V2ToV3Entry() {

	globals := &model.Globals{
		TargetDatas: helper.NewMemFile(),
	}

	globals.TableGetter = new(helper.SyncFileLoader)

	globals.SourceFileList = flag.Args()

	if err := v2tov3.Upgrade(globals); err != nil {
		log.Errorln(err)
		os.Exit(1)
		return
	}

	if *paramUpgradeOut != "" {

		for fileName, file := range globals.TargetDatas {

			fullFileName := filepath.Join(*paramUpgradeOut, "Upgraded_"+fileName)

			log.Infoln(fullFileName)

			file.Save(fullFileName)
		}

	}

}
