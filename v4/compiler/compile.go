package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/davyxu/tabtoy/v4/report"
	"github.com/pkg/errors"
)

func ParseIndexFile(g *model.Globals, indexFile string) (ret error) {

	report.Log.Debugf("Loading Index file: '%s'... ", indexFile)

	defer func() {

		switch err1 := recover().(type) {
		case *util.TableError:
			ret = err1
		case nil:
		default:
			panic(err1)
		}

	}()

	indexLoader := util.NewFileLoader(true, "")
	indexLoader.AddFile(indexFile)
	tab, err := indexLoader.GetFile(indexFile)
	if err != nil {
		return err
	}

	g.InputFiles = loadIndexTable(tab, indexFile)

	return
}

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
		if fileMeta.Mode == "Type" {
			if err := processTable(g, fileMeta); err != nil {
				return err
			}
			break
		}
	}

	for _, fileMeta := range g.InputFiles {
		if fileMeta.Mode != "Type" {
			if err := processTable(g, fileMeta); err != nil {
				return err
			}
		}
	}

	return
}

func processTable(g *model.Globals, fileMeta *model.FileMeta) error {

	file, err := g.DataFileGetter.GetFile(fileMeta.FileName)
	if err != nil {
		return errors.Wrap(err, fileMeta.FileName)
	}

	switch fileMeta.Mode {
	case "Data":
		report.Log.Infof("   (%s) %s", fileMeta.HeaderType, fileMeta.FileName)
		for _, tab := range loadDataTable(file, fileMeta, g.Types) {
			g.Datas.AddDataTable(tab)
		}
	case "KV":
		report.Log.Infof("   (%s) %s", fileMeta.HeaderType, fileMeta.FileName)
		for _, tab := range loadKVTable(file, fileMeta, g.Types) {
			tab.Mode = "KV"
			g.Datas.AddDataTable(tab)
		}
	case "Type":
		report.Log.Infof("   %s", fileMeta.FileName)
		loadTypeTable(file, fileMeta, g.Types)
	}

	return nil
}
