package v2

import (
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/printer"
	"path/filepath"
	"strings"
	"sync"
)

func getFileList(g *printer.Globals) (ret []string) {
	// 合并类型
	for _, in := range g.InputFileList {

		inputFile := in.(string)

		mergeFileList := strings.Split(inputFile, "+")

		for _, fileName := range mergeFileList {
			ret = append(ret, fileName)
		}
	}

	return
}

func cacheFile(g *printer.Globals) (fileObjByName map[string]*File) {

	var fileObjByNameGuard sync.Mutex
	fileObjByName = map[string]*File{}

	log.Infof("==========%s==========", i18n.String(i18n.Run_CacheFile))

	filelist := getFileList(g)

	var task sync.WaitGroup
	task.Add(len(filelist))

	for _, filename := range filelist {

		go func(xlsxFileName string) {

			log.Infoln(filepath.Base(xlsxFileName))
			file := NewFile(xlsxFileName)

			fileObjByNameGuard.Lock()
			fileObjByName[xlsxFileName] = file
			fileObjByNameGuard.Unlock()

			task.Done()

		}(filename)

	}

	task.Wait()

	return
}
