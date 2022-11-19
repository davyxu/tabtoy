package gen

import "github.com/davyxu/tabtoy/v4/model"

type OutputFunc func(globals *model.Globals, param string) (err error)
