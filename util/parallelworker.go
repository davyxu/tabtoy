package util

// 封装signal返回, 因为go关键字会忽略函数返回值, 所以用channel来传递结果
func task(input interface{}, callback func(interface{}) bool, signal chan bool) bool {

	result := callback(input)

	if signal != nil {
		signal <- result
		return result
	}

	return result
}

func ParallelWorker(fileList []interface{}, para bool, callback func(interface{}) bool) bool {

	// 处理多个导出文件情况

	var signal chan bool

	if para {
		signal = make(chan bool)
	}

	for _, v := range fileList {
		inputFile := v

		if signal != nil {
			go task(inputFile, callback, signal)
		} else {

			if !task(inputFile, callback, signal) {

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
