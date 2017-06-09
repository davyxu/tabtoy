package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/filter"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

func coloumnProcessor(file model.GlobalChecker, record *model.Record, fd *model.FieldDescriptor, raw string) bool {

	spliter := fd.ListSpliter()

	if fd.IsRepeated && spliter != "" {

		valueList := strings.Split(raw, spliter)

		var node *model.Node

		if fd.Type != model.FieldType_Struct {
			node = record.NewNodeByDefine(fd)
		}

		for _, v := range valueList {

			rawSingle := strings.TrimSpace(v)

			// 结构体要多添加一个节点, 处理repeated 结构体情况
			if fd.Type == model.FieldType_Struct {
				node = record.NewNodeByDefine(fd)
				node.StructRoot = true
				node = node.AddKey(fd)
			}

			if raw != "" {
				if !dataProcessor(file, fd, rawSingle, node) {
					return false
				}
			}

		}

	} else { // 普通数据/repeated单元格分多个列

		node := record.NewNodeByDefine(fd)

		// 结构体要多添加一个节点, 处理repeated 结构体情况
		if fd.Type == model.FieldType_Struct {

			node.StructRoot = true
			node = node.AddKey(fd)
		}

		if !dataProcessor(file, fd, raw, node) {
			return false
		}
	}

	return true
}

func dataProcessor(gc model.GlobalChecker, fd *model.FieldDescriptor, raw string, node *model.Node) bool {

	// 单值
	if cv, ok := filter.ConvertValue(fd, raw, gc.GlobalFileDesc(), node); !ok {
		goto ConvertError

	} else {

		// 值重复检查
		if fd.Meta.GetBool("RepeatCheck") && !gc.CheckValueRepeat(fd, cv) {
			log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueRepeated), fd.String(), cv)
			return false
		}
	}

	return true

ConvertError:

	log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueConvertError), fd.String(), raw)

	return false
}
