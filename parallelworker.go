package main

import (
	"path"
)

// 封装signal返回, 因为go关键字会忽略函数返回值, 所以用channel来传递结果
func task(input, output string, callback func(string, string) bool, signal chan bool) bool {

	result := callback(input, output)

	if signal != nil {
		signal <- result
		return result
	}

	return result
}

func parallelWorker(fileList []string, para bool, outDir string, callback func(string, string) bool) bool {

	// 处理多个导出文件情况

	var signal chan bool

	if para {
		signal = make(chan bool)
	}

	for _, v := range fileList {
		inputFile := v

		// 使用指定的导出文件夹,并更换电子表格输入文件的后缀名为pbt作为输出文件
		outputFile := path.Join(outDir, changeFileExt(inputFile, getOutputExt()))

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
