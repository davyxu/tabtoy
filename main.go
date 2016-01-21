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
	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/data"
)

var log *golog.Logger = golog.New("main")

// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
var paramPbFile = flag.String("pb", "PB", "input protobuf binary descript file, export by protoc-gen-meta plugins")

// 输入电子表格文件
var paramXlsFile = flag.String("xls", "XLS", "input excel file, use ',' splited file list by multipy files")

// 输出文件夹
var paramOut = flag.String("out", "OUT_DIR", "output directory")

// 开启调试信息
var paramDebugLevel = flag.Int("debug", 0, "show debug info")

// 并发导出,提高导出速度, 输出日志会混乱
var paramPara = flag.Bool("para", false, "parallel export by your cpu count")

func changeFileExt(filename, newExt string) string {

	file := filepath.Base(filename)

	return strings.TrimSuffix(file, path.Ext(file)) + newExt
}

func main() {

	// 关闭pbmeta的调试显示
	golog.SetLevelByString("pbmeta", "info")

	flag.Parse()

	// 调试信息挂接命令行
	data.DebuggingLevel = *paramDebugLevel

	// 加载描述文件
	fds, err := pbmeta.LoadFileDescriptorSet(*paramPbFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	// 创建描述文件池
	pool := pbmeta.NewDescriptorPool(fds)

	if !parallelWorker(*paramXlsFile, *paramPara, func(input, output string) bool {

		return export(pool, input, output)

	}) {

		halt()
		os.Exit(1)
		return
	}
}

// 封装signal返回, 因为go关键字会忽略函数返回值, 所以用channel来传递结果
func task(input, output string, callback func(string, string) bool, signal chan bool) bool {

	result := callback(input, output)

	if signal != nil {
		signal <- result
		return result
	}

	return result
}

func parallelWorker(inputFileList string, para bool, callback func(string, string) bool) bool {

	// 处理多个导出文件情况
	fileList := strings.Split(inputFileList, ",")

	var signal chan bool

	if para {
		signal = make(chan bool)
	}

	for _, v := range fileList {
		inputFile := v

		// 使用指定的导出文件夹,并更换电子表格输入文件的后缀名为pbt作为输出文件
		outputFile := path.Join(*paramOut, changeFileExt(inputFile, ".pbt"))

		if signal != nil {
			go task(inputFile, outputFile, callback, signal)
		} else {

			if !task(inputFile, outputFile, callback, signal) {

				return false
			}
		}

	}

	// 并发导出同步
	if signal != nil {
		for i := 0; i < len(fileList); i++ {
			result := <-signal
			if !result {

				return false
			}
		}
	}

	return true
}

func halt() {
	reader := bufio.NewReader(os.Stdin)

	reader.ReadLine()
}
