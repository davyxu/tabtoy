package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/filter"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

func dataProcessor(file *File, fd *model.FieldDescriptor, raw string, node *model.Node) bool {

	// 列表
	if fd.IsRepeated {

		spliter := fd.ListSpliter()

		// 使用多格子实现的repeated
		if spliter == "" {

			if _, ok := filter.ConvertValue(fd, raw, file.GlobalFD, node); !ok {
				goto ConvertError
			}

		} else {
			// 一个格子切割的repeated

			valueList := strings.Split(raw, spliter)

			for _, v := range valueList {

				if _, ok := filter.ConvertValue(fd, v, file.GlobalFD, node); !ok {
					goto ConvertError
				}
			}

		}

	} else {

		// 单值
		if cv, ok := filter.ConvertValue(fd, raw, file.GlobalFD, node); !ok {
			goto ConvertError

		} else {

			// 值重复检查
			if fd.Meta.GetBool("RepeatCheck") && !file.checkValueRepeat(fd, cv) {
				log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueRepeated), fd.String(), cv)
				return false
			}
		}

	}

	return true

ConvertError:

	log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueConvertError), fd.String(), raw)

	return false
}
