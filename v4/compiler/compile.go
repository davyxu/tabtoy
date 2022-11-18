package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/pkg/errors"
)

func Compile(g *model.Globals) (ret error) {
	defer func() {

		switch err := recover().(type) {
		case *util.TableError:
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	if g.DataFileGetter == nil {
		fileLoader := util.NewFileLoader(!g.ParaLoading, g.CacheDir)

		if g.ParaLoading {
			for _, fileMeta := range g.InputFiles {
				fileLoader.AddFile(fileMeta.FileName)
			}

			fileLoader.Commit()
		}
		g.DataFileGetter = fileLoader
	}

	for _, fileMeta := range g.InputFiles {
		file, err := g.DataFileGetter.GetFile(fileMeta.FileName)
		if err != nil {
			return errors.Wrap(err, fileMeta.FileName)
		}

		switch fileMeta.Mode {
		case "":
			for _, tab := range loadDataTable(file, fileMeta.FileName, g.Types) {
				g.Datas.AddDataTable(tab)
			}
		case "KV":
			for _, tab := range loadKVTable(file, fileMeta.FileName, g.Types) {
				g.Datas.AddDataTable(tab)
			}
		}

	}

	return
}
