package v2

import (
	"fmt"
	"strings"

	"github.com/davyxu/tabtoy/util"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

/*
	Sheet数据表单类型头

*/

const (
	// 信息所在的行
	DataSheetHeader_FieldName = 0 // 字段名(对应proto)
	DataSheetHeader_FieldType = 1 // 字段类型
	DataSheetHeader_FieldMeta = 2 // 字段特性
	DataSheetHeader_Comment   = 3 // 用户注释
	DataSheetHeader_DataBegin = 4 // 数据开始
)

type DataHeader struct {

	// 按排列的, 保留有注释掉的字段和重复的repeated列
	rawHeaderFields []*model.FieldDescriptor

	// 按排列的, 字段不重复
	headerFields []*model.FieldDescriptor

	// 按字段名分组索引字段, 字段不重复
	HeaderByName map[string]*model.FieldDescriptor
}

// 检查字段行的长度
func (self *DataHeader) ParseProtoField(index int, sheet *Sheet, localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	verticalHeader := localFD.Pragma.GetBool("Vertical")

	// 适用于配置的表格
	if verticalHeader {
		// 遍历行(从第二行开始)
		for sheet.Row = 1; ; sheet.Row++ {

			he := &DataHeaderElement{
				FieldName: sheet.GetCellData(sheet.Row, DataSheetHeader_FieldName),
				FieldType: sheet.GetCellData(sheet.Row, DataSheetHeader_FieldType),
				FieldMeta: sheet.GetCellData(sheet.Row, DataSheetHeader_FieldMeta),
				Comment:   sheet.GetCellData(sheet.Row, DataSheetHeader_Comment),
			}

			if he.FieldName == "" {
				break
			}

			if errorPos := self.addHeaderElement(he, localFD, globalFD); errorPos != -1 {
				sheet.Column = errorPos
				goto ErrorStop
			}

		}

	} else {
		// 适用于正常数据的表格

		// 遍历列
		for sheet.Column = 0; ; sheet.Column++ {

			he := &DataHeaderElement{
				FieldName: sheet.GetCellData(DataSheetHeader_FieldName, sheet.Column),
				FieldType: sheet.GetCellData(DataSheetHeader_FieldType, sheet.Column),
				FieldMeta: sheet.GetCellData(DataSheetHeader_FieldMeta, sheet.Column),
				Comment:   sheet.GetCellData(DataSheetHeader_Comment, sheet.Column),
			}

			if he.FieldName == "" {
				break
			}

			if errorPos := self.addHeaderElement(he, localFD, globalFD); errorPos != -1 {
				sheet.Row = errorPos
				goto ErrorStop
			}

		}

	}

	if len(self.rawHeaderFields) == 0 {
		return false
	}

	if index == 0 {
		// 添加第一个数据表的定义
		if !self.makeRowDescriptor(localFD, self.headerFields) {
			goto ErrorStop
		}
	}

	return true

ErrorStop:

	r, c := sheet.GetRC()

	log.Errorf("%s|%s(%s)", sheet.file.FileName, sheet.Name, util.R1C1ToA1(r, c))
	return false
}

func (self *DataHeader) RawField(index int) *model.FieldDescriptor {
	if index >= len(self.rawHeaderFields) {
		return nil
	}

	return self.rawHeaderFields[index]
}

func (self *DataHeader) RawFieldCount() int {
	return len(self.rawHeaderFields)
}

func (self *DataHeader) FieldRepeatedCount(fd *model.FieldDescriptor) (count int) {

	for _, libfd := range self.rawHeaderFields {
		if libfd == fd {
			count++
		}
	}

	return

}

func (self *DataHeader) Equal(other *DataHeader) (string, bool) {

	if len(self.headerFields) != len(other.headerFields) {
		return "field len", false
	}

	for k, v := range self.headerFields {
		if !v.Equal(other.headerFields[k]) {
			return v.Name, false
		}
	}

	return "", true
}

func (self *DataHeader) AsymmetricEqual(other *DataHeader) (string, bool) {

	for _, otherFD := range other.headerFields {

		if otherFD == nil {
			continue
		}

		if thisFD, ok := self.HeaderByName[otherFD.Name]; ok {

			if !thisFD.Equal(otherFD) {
				return otherFD.Name, false
			}
		}

	}

	return "", true
}

func (self *DataHeader) addHeaderElement(he *DataHeaderElement, localFD *model.FileDescriptor, globalFD *model.FileDescriptor) int {
	def := model.NewFieldDescriptor()
	def.Name = he.FieldName

	var errorPos int = -1

	// #开头表示注释, 跳过
	if strings.Index(he.FieldName, "#") != 0 {

		errorPos = he.Parse(def, localFD, globalFD, self.HeaderByName)
		if errorPos != -1 {
			return errorPos
		}

		// 根据字段名查找, 处理repeated字段case
		exist, ok := self.HeaderByName[def.Name]

		if ok {

			errorPos = checkSameNameElement(exist, def)
			if errorPos != -1 {
				return errorPos
			}

			def = exist

		} else {

			// 普通表头合法性检查
			errorPos = checkElement(def)
			if errorPos != -1 {
				return errorPos
			}

			self.HeaderByName[def.Name] = def
			self.headerFields = append(self.headerFields, def)
		}
	}

	// 有注释字段, 但是依然要放到这里来进行索引
	self.rawHeaderFields = append(self.rawHeaderFields, def)

	return -1
}

func (self *DataHeader) makeRowDescriptor(fileD *model.FileDescriptor, rootField []*model.FieldDescriptor) bool {

	rowType := model.NewDescriptor()
	rowType.Usage = model.DescriptorUsage_RowType
	rowType.Name = fmt.Sprintf("%sDefine", fileD.Pragma.GetString("TableName"))
	rowType.Kind = model.DescriptorKind_Struct

	// 类型已经存在, 说明是自己定义的 XXDefine, 不允许
	if _, ok := fileD.DescriptorByName[rowType.Name]; ok {
		log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_UseReservedTypeName), rowType.Name)
		return false
	}

	fileD.Add(rowType)

	// 将表格中的列添加到类型中, 方便导出
	for _, field := range rootField {

		rowType.Add(field)
	}

	return true

}

func newDataHeadSheet() *DataHeader {

	return &DataHeader{
		HeaderByName: make(map[string]*model.FieldDescriptor),
	}
}
