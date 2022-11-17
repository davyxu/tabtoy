package report

import (
	"github.com/davyxu/golog"
)

var Log = golog.New("tabtoy_v4")

func init() {
	Log.SetParts()
}
