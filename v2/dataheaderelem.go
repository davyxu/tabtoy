package v2

import (
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
	"strings"
)

type DataHeaderElement struct {
	FieldName string
	FieldType string
	FieldMeta string
	Comment   string
}

func checkElement(def *model.FieldDescriptor) int {
	// 普通表头合法性检查

	// 结构体单元格不能进行切分
	//if def.Type == model.FieldType_Struct && def.ListSpliter() != "" {
	//
	//	log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_StructCellCannotSplit),
	//		def.Name)
	//	return DataSheetHeader_FieldMeta
	//}

	return -1
}

func checkSameNameElement(exist, def *model.FieldDescriptor) int {
	// 多个同名字段只允许repeated方式的字段
	if !exist.IsRepeated {
		log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_DuplicateFieldName), def.Name)
		return DataSheetHeader_FieldName
	}

	// 多个repeated描述类型不一致
	if exist.Type != def.Type {

		log.Errorf("%s '%s' '%s' '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
			def.Name,
			model.FieldTypeToString(exist.Type),
			model.FieldTypeToString(def.Type))

		return DataSheetHeader_FieldType
	}

	// 多个repeated描述内建类型不一致
	if exist.Complex != def.Complex {

		log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
			def.Name)

		return DataSheetHeader_FieldType
	}

	// 多个repeated描述的meta不一致
	if exist.Meta.String() != def.Meta.String() {

		log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldMetaNotSameInMultiColumn),
			def.Name)

		return DataSheetHeader_FieldMeta
	}

	return -1
}

func (self *DataHeaderElement) Parse(def *model.FieldDescriptor, localFD *model.FileDescriptor, globalFD *model.FileDescriptor, headerByName map[string]*model.FieldDescriptor) int {

	// ====================解析类型====================

	testFileD := localFD

	for {

		if def.ParseType(testFileD, self.FieldType) {
			break
		}

		if testFileD == localFD {
			testFileD = globalFD
			continue
		}

		break
	}

	// 依然找不到, 报错
	if def.Type == model.FieldType_None {
		log.Errorf("%s, '%s' (%s) raw: %s", i18n.String(i18n.DataHeader_TypeNotFound), def.Name, model.FieldTypeToString(def.Type), self.FieldType)
		return DataSheetHeader_FieldType
	}

	// ====================解析特性====================
	if err := def.Meta.Parse(self.FieldMeta); err != nil {
		log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_MetaParseFailed), err)
		return DataSheetHeader_FieldMeta
	}

	def.Comment = strings.Replace(self.Comment, "\n", " ", -1)

	// 根据字段名查找, 处理repeated字段case
	exist, ok := headerByName[def.Name]

	if ok {

		// 多个同名字段只允许repeated方式的字段
		if !exist.IsRepeated {
			log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_DuplicateFieldName), def.Name)
			return DataSheetHeader_FieldName
		}

		// 多个repeated描述类型不一致
		if exist.Type != def.Type {

			log.Errorf("%s '%s' '%s' '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
				def.Name,
				model.FieldTypeToString(exist.Type),
				model.FieldTypeToString(def.Type))

			return DataSheetHeader_FieldType
		}

		// 多个repeated描述内建类型不一致
		if exist.Complex != def.Complex {

			log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
				def.Name)

			return DataSheetHeader_FieldType
		}

		// 多个repeated描述的meta不一致
		if exist.Meta.String() != def.Meta.String() {

			log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldMetaNotSameInMultiColumn),
				def.Name)

			return DataSheetHeader_FieldMeta
		}

		def = exist
	}

	return -1
}
