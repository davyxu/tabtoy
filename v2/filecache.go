package v2

import (
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/printer"
	"os"
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

	var cachedir string
	if g.UseCache {
		cachedir = g.CacheDir

		os.Mkdir(g.CacheDir, 0666)
	}

	writeOK := func(xlsxFileName string) {
		file, fromCache := NewFile(xlsxFileName, cachedir)

		fileObjByNameGuard.Lock()
		fileObjByName[xlsxFileName] = file
		fileObjByNameGuard.Unlock()
		if fromCache {
			log.Infof("%s [Cache]", filepath.Base(xlsxFileName))
		} else {
			log.Infof("%s", filepath.Base(xlsxFileName))
		}

	}

	// 这里加速效果良好, 默认开启
	var task sync.WaitGroup
	task.Add(len(filelist))

	for _, filename := range filelist {

		go func(xlsxFileName string) {

			writeOK(xlsxFileName)
			task.Done()

		}(filename)

	}

	task.Wait()

	// 调试用
	//for _, filename := range filelist {
	//
	//	writeOK(filename)
	//}

	return
}
