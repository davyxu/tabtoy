package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davyxu/goarg"
	"github.com/davyxu/golog"
	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/data"
)

var log *golog.Logger = golog.New("main")

func changeFileExt(filename, newExt string) string {

	file := filepath.Base(filename)

	return strings.TrimSuffix(file, path.Ext(file)) + newExt
}

func printHelp() {
	fmt.Println(`Usage: tabtoy [OPTION] --pb=FILE XLS_FILES
  --pb=FILE			Input protobuf binary descript file, export by protoc-gen-meta plugins.
  --out=DIR			Output directory.
  --para=n			The number of export that can be run in parallel, default to number of CPUs.
  --debug=LEVEL		Show debug info, level range( 0~5 ).
  --version			Show version.
	`)
}

func main() {

	cmdline := goarg.NewCommandLineParser()

	cmdline.Parse()

	// 关闭pbmeta的调试显示
	golog.SetLevelByString("pbmeta", "info")

	if cmdline.ArgCount() == 0 {
		printHelp()
		return
	}

	// 帮助
	if cmdline.NamedExists("--help") || cmdline.NamedExists("-h") {
		printHelp()
		return
	}

	// 版本
	if cmdline.NamedExists("--version") {
		fmt.Println("tabtoy 0.1.0")
		return
	}

	// 调试信息挂接命令行
	data.DebuggingLevel = cmdline.NamedAsInteger("--debug", 0)

	// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
	pbFile := cmdline.NamedAsString("--pb", "")

	// 加载描述文件
	fds, err := pbmeta.LoadFileDescriptorSet(pbFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 创建描述文件池
	pool := pbmeta.NewDescriptorPool(fds)

	xlsFileList := cmdline.FreeValues()
	outDir := cmdline.NamedAsString("--out", ".")

	var paraCount int

	// 有参数时, 默认为CPU个数,除非指定, 没参数时不开并发
	if cmdline.NamedExists("--para") {
		paraCount = cmdline.NamedAsInteger("--para", runtime.NumCPU())
	}

	if !parallelWorker(xlsFileList, paraCount, outDir, func(input, output string) bool {

		return export(pool, input, output)

	}) {

		halt()
		os.Exit(1)
		return
	}
}

func halt() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadLine()
}
