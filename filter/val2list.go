package filter

import (
	"strings"

	"github.com/davyxu/tabtoy/proto/tool"
)

// 分割字符串并调用回调
func Value2List(meta *tool.FieldMeta, value string, callback func(string)) bool {

	if meta == nil {
		return false
	}

	if meta.String2ListSpliter == "" {
		return false
	}

	valueList := strings.Split(value, meta.String2ListSpliter)

	for _, v := range valueList {
		callback(v)
	}

	return true
}
