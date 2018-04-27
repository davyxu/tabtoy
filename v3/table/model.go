package table

import (
	"encoding/json"
)

var (
	BuiltinConfig Config // 内嵌数据
)

func init() {
	err := json.Unmarshal([]byte(builtinJson), &BuiltinConfig)
	if err != nil {
		panic(err)
	}
}
