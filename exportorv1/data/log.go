package data

import (
	"github.com/davyxu/golog"
)

var log *golog.Logger = golog.New("data")

// 级别越高,可见信息越多(0~5)
var DebuggingLevel int
