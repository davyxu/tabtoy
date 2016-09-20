package exportorv2

import (
	"fmt"
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
	"github.com/golang/protobuf/proto"
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

// 检查字段行的长度
func (self *DataHeader) ParseProtoField(sheet *Sheet, localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	var def *model.FieldDescriptor

	// 遍历列
	for sheet.Column = 0; ; sheet.Column++ {

		def = new(model.FieldDescriptor)

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
				log.Errorf("field header type not found: '%s' (%s) raw: %s", def.Name, model.FieldTypeToString(def.Type), rawFieldType)
				goto ErrorStop
			}

			// ====================解析特性====================
			metaString := sheet.GetCellData(DataSheetRow_FieldMeta, sheet.Column)

			if err := proto.UnmarshalText(metaString, &def.Meta); err != nil {
				sheet.Row = DataSheetRow_FieldMeta
				log.Errorln("parse field header failed", err)
				goto ErrorStop
			}

			def.Comment = sheet.GetCellData(DataSheetRow_Comment, sheet.Column)

			// 根据字段名查找, 处理repeated字段case
			exist, ok := self.HeaderByName[def.Name]

			if ok {

				// 多个同名字段只允许repeated方式的字段
				if !exist.IsRepeated {
					sheet.Row = DataSheetRow_FieldName
					log.Errorf("duplicate field header: '%s'", def.Name)
					goto ErrorStop
				}

				// 多个repeated描述类型不一致
				if exist.Type != def.Type {
					sheet.Row = DataSheetRow_FieldType
					log.Errorf("repeated field type diff in multi column: '%s', prev: '%s', found: '%s'",
						def.Name,
						model.FieldTypeToString(exist.Type),
						model.FieldTypeToString(def.Type))

					goto ErrorStop
				}

				// 多个repeated描述内建类型不一致
				if exist.Complex != def.Complex {
					sheet.Row = DataSheetRow_FieldType
					log.Errorf("repeated field build type diff in multi column: '%s'",
						def.Name)

					goto ErrorStop
				}

				// 多个repeated描述的meta不一致
				if proto.CompactTextString(&exist.Meta) != proto.CompactTextString(&def.Meta) {
					sheet.Row = DataSheetRow_FieldMeta
					log.Errorf("repeated field meta diff in multi column: '%s'",
						def.Name)

					goto ErrorStop
				}

				def = exist

			} else {
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

	// 添加一次行结构
	self.makeRowDescriptor(localFD, self.headerFields)

	return true

ErrorStop:

	r, c := sheet.GetRC()

	log.Errorf("%s|%s(%s)", sheet.file.FileName, sheet.Name, util.ConvR1C1toA1(r, c))
	return false
}

func (self *DataHeader) makeRowDescriptor(fileD *model.FileDescriptor, rootField []*model.FieldDescriptor) {

	rowType := model.NewDescriptor()
	rowType.Usage = model.DescriptorUsage_RowType
	rowType.Name = fmt.Sprintf("%sDefine", fileD.Pragma.TableName)
	rowType.Kind = model.DescriptorKind_Struct

	// 有就不添加
	if _, ok := fileD.DescriptorByName[rowType.Name]; ok {
		return
	}

	fileD.Add(rowType)

	// 将表格中的列添加到类型中, 方便导出
	for _, field := range rootField {

		rowType.Add(field)
	}

}

func newDataHeadSheet() *DataHeader {

	return &DataHeader{
		HeaderByName: make(map[string]*model.FieldDescriptor),
	}
}
