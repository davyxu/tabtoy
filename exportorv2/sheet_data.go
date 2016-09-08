package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/filter"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
	"github.com/golang/protobuf/proto"
)

const (
	// 信息所在的行
	DataSheetRow_FieldName = 0 // 字段名(对应proto)
	DataSheetRow_FieldType = 1 // 字段类型
	DataSheetRow_FieldMeta = 2 // 字段特性
	DataSheetRow_Comment   = 3 // 用户注释
	DataSheetRow_DataBegin = 4 // 数据开始
)

type DataSheet struct {
	*Sheet

	// 按排列的, 合并repeated描述字段
	headerFields []*model.FieldDefine

	// 按字段名分组索引字段
	headerByName map[string]*model.FieldDefine
}

func (self *DataSheet) FetchFieldDefine(index int) *model.FieldDefine {
	if index >= len(self.headerFields) {
		return nil
	}

	return self.headerFields[index]
}

func (self *DataSheet) Valid() bool {
	return self.GetCellData(0, 0) != ""
}

// 检查字段行的长度
func (self *DataSheet) ParseProtoField(tts *model.BuildInTypeSet) bool {

	var def *model.FieldDefine

	// 遍历列
	for self.Column = 0; ; self.Column++ {

		def = new(model.FieldDefine)

		// ====================解析字段====================
		def.Name = self.GetCellData(DataSheetRow_FieldName, self.Column)
		if def.Name == "" {
			break
		}

		var colIgnore bool
		// #开头表示注释, 跳过
		if strings.Index(def.Name, "#") == 0 {
			colIgnore = true
		}

		// ====================解析类型====================
		def.ParseType(tts, self.GetCellData(DataSheetRow_FieldType, self.Column))

		// 依然找不到, 报错
		if !colIgnore && def.Type == model.FieldType_None {
			self.Row = DataSheetRow_FieldType
			log.Errorf("field header type not found: %s  %s", def.Name, model.FieldTypeToString(def.Type))
			goto ErrorStop
		}

		// ====================解析特性====================
		metaString := self.GetCellData(DataSheetRow_FieldMeta, self.Column)

		if err := proto.UnmarshalText(metaString, &def.Meta); err != nil {
			self.Row = DataSheetRow_FieldMeta
			log.Errorln("parse field header failed", err)
			goto ErrorStop
		}

		def.Comment = self.GetCellData(DataSheetRow_Comment, self.Column)

		// 根据字段名查找, 处理repeated字段case
		exist, ok := self.headerByName[def.Name]

		if ok {

			// 多个同名字段只允许repeated方式的字段
			if !exist.IsRepeated {
				self.Row = DataSheetRow_FieldName
				log.Errorf("duplicate field header: %s", def.Name)
				goto ErrorStop
			}

			// 多个repeated描述类型不一致
			if exist.Type != def.Type {
				self.Row = DataSheetRow_FieldType
				log.Errorf("repeated field type diff in multi column: %s, prev: %s, found: %s",
					def.Name,
					model.FieldTypeToString(exist.Type),
					model.FieldTypeToString(def.Type))

				goto ErrorStop
			}

			// 多个repeated描述内建类型不一致
			if exist.BuildInType != def.BuildInType {
				self.Row = DataSheetRow_FieldType
				log.Errorf("repeated field build type diff in multi column: %s",
					def.Name)

				goto ErrorStop
			}

			// 多个repeated描述的meta不一致
			if proto.CompactTextString(&exist.Meta) != proto.CompactTextString(&def.Meta) {
				self.Row = DataSheetRow_FieldMeta
				log.Errorf("repeated field meta diff in multi column: %s",
					def.Name)

				goto ErrorStop
			}

			def = exist

		} else {
			self.headerByName[def.Name] = def
		}

		self.headerFields = append(self.headerFields, def)
	}

	return len(self.headerFields) > 0

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

func dataProcessor(file *File, fieldDef *model.FieldDefine, rawValue string, node *model.Node) bool {

	//	if node.Define.Name == "Type" {
	//		a := 1
	//		a++
	//	}

	// 列表
	if fieldDef.IsRepeated {

		spliter := fieldDef.ListSpliter()

		// 使用多格子实现的repeated
		if spliter == "" {

			if _, ok := filter.ConvertValue(fieldDef, rawValue, file.TypeSet, node); !ok {
				goto ConvertError
			}

		} else {
			// 一个格子切割的repeated

			valueList := strings.Split(rawValue, spliter)

			for _, v := range valueList {

				if _, ok := filter.ConvertValue(fieldDef, v, file.TypeSet, node); !ok {
					goto ConvertError
				}
			}

		}

	} else {

		// 单值
		if cv, ok := filter.ConvertValue(fieldDef, rawValue, file.TypeSet, node); !ok {
			goto ConvertError

		} else {

			// 值重复检查
			if fieldDef.Meta.RepeatCheck && !file.checkValueRepeat(fieldDef, cv) {
				log.Errorf("found repeat value, %s raw: '%s'", fieldDef.String(), cv)
				return false
			}
		}

	}

	return true

ConvertError:

	log.Errorf("value convert error, %s raw: '%s'", fieldDef.String(), rawValue)

	return false
}

func (self *DataSheet) Export(file *File, tab *model.Table) bool {

	// 检查引导头
	if !self.ParseProtoField(file.TypeSet) {
		return true
	}

	// 是否继续读行
	var readingLine bool = true

	// 遍历每一行
	for self.Row = DataSheetRow_DataBegin; readingLine; self.Row++ {

		// 第一列是空的，结束
		if self.GetCellData(self.Row, 0) == "" {
			break
		}

		record := model.NewRecord()

		// 遍历每一列
		for self.Column = 0; self.Column < len(self.headerFields); self.Column++ {

			fieldDef := self.FetchFieldDefine(self.Column)

			// 数据大于列头时, 结束这个列
			if fieldDef == nil {
				break
			}

			// #开头表示注释, 跳过
			if strings.Index(fieldDef.Name, "#") == 0 {
				continue
			}

			rawValue := self.GetCellData(self.Row, self.Column)

			// repeated的, 没有填充的, 直接跳过, 不生成数据
			if rawValue == "" && fieldDef.IsRepeated {
				continue
			}

			node := record.NewNodeByDefine(fieldDef)

			// 结构体要多添加一个节点, 处理repeated 结构体情况
			if fieldDef.Type == model.FieldType_Struct {

				node.StructRoot = true
				node = node.AddKey(fieldDef)

			}

			if !dataProcessor(file, fieldDef, rawValue, node) {
				goto ErrorStop
			}

		}

		tab.Add(record)

	}

	return true

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

func newDataSheet(sheet *Sheet) *DataSheet {

	return &DataSheet{
		Sheet:        sheet,
		headerByName: make(map[string]*model.FieldDefine),
	}
}
