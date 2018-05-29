package report

import (
	"github.com/davyxu/golog"
)

var Log = golog.New("tabtoy2")

func init() {
	Log.SetParts()
}
