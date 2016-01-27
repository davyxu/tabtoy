package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/davyxu/golog"
	"github.com/davyxu/tabtoy/data"
)

var log *golog.Logger = golog.New("main")

// 开启调试信息
var paramDebugLevel = flag.Int("debug", 0, "show debug info")

// 并发导出,提高导出速度, 输出日志会混乱
var paramPara = flag.Bool("para", false, "parallel export by your cpu count")

// 显示版本号
var paramVersion = flag.Bool("version", false, "Show version")

// 工作模式
var paramMode = flag.String("mode", "", "mode: xls2pbt, syncheader")

func changeFileExt(filename, newExt string) string {

	file := filepath.Base(filename)

	return strings.TrimSuffix(file, path.Ext(file)) + newExt
}

func main() {

	flag.Parse()

	// 关闭pbmeta的调试显示
	golog.SetLevelByString("pbmeta", "info")

	// 版本
	if *paramVersion {
		fmt.Println("tabtoy 0.1.1")
		return
	}

	// 调试信息挂接命令行
	data.DebuggingLevel = *paramDebugLevel

	switch *paramMode {
	case "xls2pbt":
		if !runXls2PbtMode() {
			os.Exit(1)
			return
		}

		//	case "syncheader":
		//		if !runSyncHeaderMode() {
		//			os.Exit(1)
		//			return
		//		}

	default:
		fmt.Println("--mode not specify")
		os.Exit(1)
		return
	}

}

func halt() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadLine()
}
