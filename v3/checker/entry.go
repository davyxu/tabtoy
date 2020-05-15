package checker

import "github.com/davyxu/tabtoy/v3/model"

func CheckData(globals *model.Globals) {

	checkEnumValue(globals)
	checkRepeat(globals)
	checkDataType(globals)
}
