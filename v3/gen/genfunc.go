package gen

import "github.com/davyxu/tabtoy/v3/model"

type GenFunc func(globals *model.Globals) (data []byte, err error)
