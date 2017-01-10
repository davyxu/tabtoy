package exportorv2

import (
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

/*
	@Types表解析

*/

const (
	// 信息所在的行
	TypeSheetRow_Pragma    = 0 // 配置
	TypeSheetRow_FieldDesc = 1 // 类型描述
	TypeSheetRow_Comment   = 2 // 字段名(对应proto)
	TypeSheetRow_DataBegin = 3 // 数据开始
)

var typeHeader = map[string]int{
	"ObjectType": 0,
	"FieldName":  1,
	"FieldType":  2,
	"Value":      3,
	"Comment":    4,
	"Meta":       5,
	"Alias":      6,
	"Default":    7,
}

type TypeSheet struct {
	*Sheet
}

func (self *TypeSheet) parseTable(root *typeModelRoot) bool {

	var readingLine bool = true

	root.pragma = self.GetCellData(TypeSheetRow_Pragma, 0)

	// 读行
	for row := TypeSheetRow_DataBegin; readingLine; row++ {

		tm := newTypeModel()
		tm.row = row

		// 读列
		for col := 0; ; col++ {

			// 头的类型
			typeDeclare := self.GetCellData(TypeSheetRow_FieldDesc, col)

			// 头已经读完
			if typeDeclare == "" {
				break
			}

			if _, ok := typeHeader[typeDeclare]; !ok {
				self.Row = TypeSheetRow_FieldDesc
				self.Column = col
				log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_UnexpectedTypeHeader), typeDeclare)
				return false
			}

			// 值
			typeValue := self.GetCellData(row, col)

			// 类型空表示停止解析
			if typeDeclare == "ObjectType" && typeValue == "" {
				readingLine = false
				break
			}

			tm.colData[typeDeclare] = &typeCell{
				value: typeValue,
				col:   col,
			}

		}

		if len(tm.colData) > 0 {
			root.models = append(root.models, tm)
		}

	}

	return true
}

// 解析所有的类型及数据
func (self *TypeSheet) Parse(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	var root typeModelRoot

	if !self.parseTable(&root) {
		goto ErrorStop
	}

	if !root.ParsePragma(localFD) {
		self.Row = TypeSheetRow_Pragma
		self.Column = 0
		log.Errorf("%s", i18n.String(i18n.TypeSheet_PackageIsEmpty))
		goto ErrorStop
	}

	if !root.ParseData(localFD, globalFD) {
		self.Row = root.Row
		self.Column = root.Col
		goto ErrorStop
	}

	if !root.SolveUnknownModel(localFD, globalFD) {
		self.Row = root.Row
		self.Column = root.Col
		goto ErrorStop
	}

	return self.checkProtobufCompatibility(localFD)

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

// 检查protobuf兼容性
func (self *TypeSheet) checkProtobufCompatibility(fileD *model.FileDescriptor) bool {

	for _, bt := range fileD.Descriptors {
		if bt.Kind == model.DescriptorKind_Enum {

			// proto3 需要枚举有0值
			if _, ok := bt.FieldByNumber[0]; !ok {
				log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_FirstEnumValueShouldBeZero), bt.Name)
				return false
			}
		}
	}

	return true
}

func newTypeSheet(sheet *Sheet) *TypeSheet {
	return &TypeSheet{
		Sheet: sheet,
	}
}
