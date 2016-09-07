package exportorv2

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
	"github.com/golang/protobuf/proto"
)

const (
	// 信息所在的行
	TypeSheetRow_Comment   = 0 // 字段名(对应proto)
	TypeSheetRow_DataBegin = 1 // 数据开始
)

const (
	// 信息所在列
	TypeSheetCol_Type      = 0 // 类型
	TypeSheetCol_FieldName = 1 // 字段名
	TypeSheetCol_Value     = 2 // 值
	TypeSheetCol_Meta      = 3 // 特性
)

type TypeSheet struct {
	*Sheet

	*model.BuildInTypeSet
}

func (self *TypeSheet) Parse() bool {

	// 是否继续读行
	var readingLine bool = true

	var td *model.BuildInType

	// 遍历每一行
	for self.Row = TypeSheetRow_DataBegin; readingLine; self.Row++ {

		// 第一列是空的，结束
		if self.GetCellData(self.Row, TypeSheetCol_Type) == "" {
			break
		}

		var fd model.BuildInTypeField

		// 解析枚举类型
		rawTypeName := self.GetCellData(self.Row, TypeSheetCol_Type)

		existType, ok := self.BuildInTypeSet.TypeByName[rawTypeName]

		if ok {

			td = existType

		} else {

			td = model.NewBuildInType()
			td.Name = rawTypeName
			self.BuildInTypeSet.Add(td)
		}

		// 解析字段名
		fd.Name = self.GetCellData(self.Row, TypeSheetCol_FieldName)

		rawValue := self.GetCellData(self.Row, TypeSheetCol_Value)

		var kind model.BuildInTypeKind

		// 非空值是枚举
		if rawValue != "" {

			// 解析枚举值
			if v, err := strconv.Atoi(rawValue); err == nil {
				fd.Value = int32(v)
			} else {
				self.Column = TypeSheetCol_Value
				log.Errorln("parse type value failed:", err)
				goto ErrorStop
			}
			kind = model.BuildInTypeKind_Enum
		} else {
			kind = model.BuildInTypeKind_Struct
		}

		if td.Kind == model.BuildInTypeKind_None {
			td.Kind = kind
			// 一些字段有填值, 一些没填值
		} else if td.Kind != kind {
			self.Column = TypeSheetCol_Value
			log.Errorln("buildin kind shold be same", td.Kind, kind)
			goto ErrorStop
		}

		// 解析特性
		metaString := self.GetCellData(self.Row, TypeSheetCol_Meta)

		if err := proto.UnmarshalText(metaString, &fd.Meta); err != nil {
			log.Errorln("parse field header failed", err)
			return false
		}

		td.Add(&fd)

	}

	return true

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

func newTypeSheet(sheet *Sheet) *TypeSheet {
	return &TypeSheet{
		Sheet:          sheet,
		BuildInTypeSet: model.NewBuildInTypeSet(),
	}
}
