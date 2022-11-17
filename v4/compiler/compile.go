package compiler

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v4/model"
	"github.com/pkg/errors"
)

func Compile(cp *model.Compiler) (ret error) {
	defer func() {

		switch err := recover().(type) {
		case *util.TableError:
			ret = err
		case nil:
		default:
			panic(err)
		}

	}()

	if cp.DataFileGetter == nil {
		fileLoader := util.NewFileLoader(!cp.ParaLoading, cp.CacheDir)

		if cp.ParaLoading {
			for _, fileMeta := range cp.InputFiles {
				fileLoader.AddFile(fileMeta.FileName)
			}

			fileLoader.Commit()
		}
		cp.DataFileGetter = fileLoader
	}

	for _, fileMeta := range cp.InputFiles {
		file, err := cp.DataFileGetter.GetFile(fileMeta.FileName)
		if err != nil {
			return errors.Wrap(err, fileMeta.FileName)
		}

		switch fileMeta.Mode {
		case "":
			for _, tab := range loadDataTable(file, fileMeta.FileName, cp.Types) {
				cp.Datas.AddDataTable(tab)
			}
		case "KV":
			for _, tab := range loadKVTable(file, fileMeta.FileName, cp.Types) {
				cp.Datas.AddDataTable(tab)
			}
		}

	}

	return
}
