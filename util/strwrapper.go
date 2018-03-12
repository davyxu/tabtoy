package util

import (
	"fmt"
)

func StringEscape(s string) string {

	b := make([]byte, 0)

	var index int

	// 表中直接使用换行会干扰最终合并文件格式, 所以转成\n,由pbt文本解析层转回去
	for index < len(s) {
		c := s[index]

		switch c {
		case '"':
			b = append(b, '\\')
			b = append(b, '"')
		case '\n':
			b = append(b, '\\')
			b = append(b, 'n')
		case '\r':
			b = append(b, '\\')
			b = append(b, 'r')
		case '\\':
			b = append(b, '\\')
			b = append(b, c)
		default:
			b = append(b, c)
		}

		index++

	}

	return fmt.Sprintf("\"%s\"", string(b))

}
