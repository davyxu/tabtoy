package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
	"github.com/golang/protobuf/proto"
)

type DataHeader struct {

	// 按排列的, 合并repeated描述字段
	headerFields []*model.FieldDefine

	// 按字段名分组索引字段
	HeaderByName map[string]*model.FieldDefine
}

func (self *DataHeader) FetchFieldDefine(index int) *model.FieldDefine {
	if index >= len(self.headerFields) {
		return nil
	}

	return self.headerFields[index]
}

// 检查字段行的长度
func (self *DataHeader) ParseProtoField(sheet *Sheet, tts *model.BuildInTypeSet) bool {

	var def *model.FieldDefine

	// 遍历列
	for sheet.Column = 0; ; sheet.Column++ {

		def = new(model.FieldDefine)

		// ====================解析字段====================
		def.Name = sheet.GetCellData(DataSheetRow_FieldName, sheet.Column)
		if def.Name == "" {
			break
		}

		var colIgnore bool
		// #开头表示注释, 跳过
		if strings.Index(def.Name, "#") == 0 {
			colIgnore = true
		}

		// ====================解析类型====================
		def.ParseType(tts, sheet.GetCellData(DataSheetRow_FieldType, sheet.Column))

		// 依然找不到, 报错
		if !colIgnore && def.Type == model.FieldType_None {
			sheet.Row = DataSheetRow_FieldType
			log.Errorf("field header type not found: %s  %s", def.Name, model.FieldTypeToString(def.Type))
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
				log.Errorf("duplicate field header: %s", def.Name)
				goto ErrorStop
			}

			// 多个repeated描述类型不一致
			if exist.Type != def.Type {
				sheet.Row = DataSheetRow_FieldType
				log.Errorf("repeated field type diff in multi column: %s, prev: %s, found: %s",
					def.Name,
					model.FieldTypeToString(exist.Type),
					model.FieldTypeToString(def.Type))

				goto ErrorStop
			}

			// 多个repeated描述内建类型不一致
			if exist.BuildInType != def.BuildInType {
				sheet.Row = DataSheetRow_FieldType
				log.Errorf("repeated field build type diff in multi column: %s",
					def.Name)

				goto ErrorStop
			}

			// 多个repeated描述的meta不一致
			if proto.CompactTextString(&exist.Meta) != proto.CompactTextString(&def.Meta) {
				sheet.Row = DataSheetRow_FieldMeta
				log.Errorf("repeated field meta diff in multi column: %s",
					def.Name)

				goto ErrorStop
			}

			def = exist

		} else {
			self.HeaderByName[def.Name] = def
		}

		self.headerFields = append(self.headerFields, def)
	}

	return len(self.headerFields) > 0

ErrorStop:

	r, c := sheet.GetRC()

	log.Errorf("%s|%s(%s)", sheet.file.FileName, sheet.Name, util.ConvR1C1toA1(r, c))
	return false
}

func newDataHeadSheet() *DataHeader {

	return &DataHeader{
		HeaderByName: make(map[string]*model.FieldDefine),
	}
}
