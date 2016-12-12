package exportorv2

import (
	"fmt"
	"strings"

	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

type DataHeader struct {

	// 按排列的, 保留有注释掉的字段和重复的repeated列
	rawHeaderFields []*model.FieldDescriptor

	// 按排列的, 字段不重复
	headerFields []*model.FieldDescriptor

	// 按字段名分组索引字段, 字段不重复
	HeaderByName map[string]*model.FieldDescriptor
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

func (self *DataHeader) Equal(other *DataHeader) bool {

	if len(self.headerFields) != len(other.headerFields) {
		return false
	}

	for k, v := range self.headerFields {
		if !v.Equal(other.headerFields[k]) {
			return false
		}
	}

	return true
}

// 检查字段行的长度
func (self *DataHeader) ParseProtoField(index int, sheet *Sheet, localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	var def *model.FieldDescriptor

	// 遍历列
	for sheet.Column = 0; ; sheet.Column++ {

		def = model.NewFieldDescriptor()

		// ====================解析字段====================
		def.Name = sheet.GetCellData(DataSheetRow_FieldName, sheet.Column)
		if def.Name == "" {
			break
		}

		// #开头表示注释, 跳过
		if strings.Index(def.Name, "#") != 0 {

			// ====================解析类型====================

			testFileD := localFD

			rawFieldType := sheet.GetCellData(DataSheetRow_FieldType, sheet.Column)

			for {
				if def.ParseType(testFileD, rawFieldType) {
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
				sheet.Row = DataSheetRow_FieldType
				log.Errorf("%s, '%s' (%s) raw: %s", i18n.String(i18n.DataHeader_TypeNotFound), def.Name, model.FieldTypeToString(def.Type), rawFieldType)
				goto ErrorStop
			}

			// ====================解析特性====================
			metaString := sheet.GetCellData(DataSheetRow_FieldMeta, sheet.Column)

			def.Meta = golexer.NewKVPair()

			if err := def.Meta.Parse(metaString); err != nil {
				sheet.Row = DataSheetRow_FieldMeta
				log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_MetaParseFailed), err)
				goto ErrorStop
			}

			def.Comment = sheet.GetCellData(DataSheetRow_Comment, sheet.Column)

			// 根据字段名查找, 处理repeated字段case
			exist, ok := self.HeaderByName[def.Name]

			if ok {

				// 多个同名字段只允许repeated方式的字段
				if !exist.IsRepeated {
					sheet.Row = DataSheetRow_FieldName
					log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_DuplicateFieldName), def.Name)
					goto ErrorStop
				}

				// 多个repeated描述类型不一致
				if exist.Type != def.Type {
					sheet.Row = DataSheetRow_FieldType

					log.Errorf("%s '%s' '%s' '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
						def.Name,
						model.FieldTypeToString(exist.Type),
						model.FieldTypeToString(def.Type))

					goto ErrorStop
				}

				// 多个repeated描述内建类型不一致
				if exist.Complex != def.Complex {
					sheet.Row = DataSheetRow_FieldType

					log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldTypeNotSameInMultiColumn),
						def.Name)

					goto ErrorStop
				}

				// 多个repeated描述的meta不一致
				if exist.Meta.String() != def.Meta.String() {
					sheet.Row = DataSheetRow_FieldMeta

					log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_RepeatedFieldMetaNotSameInMultiColumn),
						def.Name)

					goto ErrorStop
				}

				def = exist

			} else {

				// 普通表头合法性检查

				// 结构体单元格不能进行切分
				if def.Type == model.FieldType_Struct && def.ListSpliter() != "" {
					sheet.Row = DataSheetRow_FieldMeta
					log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_StructCellCannotSplit),
						def.Name)
					goto ErrorStop
				}

				self.HeaderByName[def.Name] = def
				self.headerFields = append(self.headerFields, def)
			}
		}

		// 有注释字段, 但是依然要放到这里来进行索引
		self.rawHeaderFields = append(self.rawHeaderFields, def)
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

	log.Errorf("%s|%s(%s)", sheet.file.FileName, sheet.Name, util.ConvR1C1toA1(r, c))
	return false
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
