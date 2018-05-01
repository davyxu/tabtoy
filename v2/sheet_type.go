package v2

import (
	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
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

func (self *TypeSheet) detectMaxTypeCol() int {
	for col := 0; col < 100; col++ {

		// 头的类型
		typeDeclare := self.GetCellData(TypeSheetRow_FieldDesc, col)

		// 头已经读完
		if typeDeclare == "" {
			return col
		}
	}

	return 0
}

func (self *TypeSheet) parseTable(root *typeModelRoot) bool {

	var readingLine bool = true

	root.pragma = self.GetCellData(TypeSheetRow_Pragma, 0)

	maxCol := self.detectMaxTypeCol()

	var meetEmptyLine bool

	var warningAfterEmptyLineDataOnce bool

	// 读行
	for row := TypeSheetRow_DataBegin; readingLine; row++ {

		tm := newTypeModel()
		tm.row = row

		// 整行都是空的
		if self.IsFullRowEmpty(row, maxCol) {

			// 再次碰空行, 表示确实是空的
			if meetEmptyLine {
				break

			} else {
				meetEmptyLine = true
			}

			continue

		} else {

			//已经碰过空行, 这里又碰到数据, 说明有人为隔出的空行, 做warning提醒, 防止数据没导出
			if meetEmptyLine && !warningAfterEmptyLineDataOnce {
				log.Errorf("%s %s|%s(%s)", i18n.String(i18n.TypeSheet_RowDataSplitedByEmptyLine), self.file.FileName, self.Name, util.R1C1ToA1(row, 1))

				warningAfterEmptyLineDataOnce = true
			}

		}

		// 读列
		for col := 0; col < maxCol; col++ {

			// 头的类型
			typeDeclare := self.GetCellData(TypeSheetRow_FieldDesc, col)

			if _, ok := typeHeader[typeDeclare]; !ok {
				self.Row = TypeSheetRow_FieldDesc
				self.Column = col
				log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_UnexpectedTypeHeader), typeDeclare)
				return false
			}

			// 值
			typeValue := self.GetCellData(row, col)

			if typeDeclare == "ObjectType" && typeValue == "" {
				self.Row = row
				self.Column = col
				log.Errorf("%s", i18n.String(i18n.TypeSheet_ObjectNameEmpty))
				return false
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

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.R1C1ToA1(r, c))
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
