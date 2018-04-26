package table

import (
	"encoding/json"
)

var (
	config Config // 内嵌数据
)

func init() {
	err := json.Unmarshal([]byte(builtinJson), &config)
	if err != nil {
		panic(err)
	}
}
