package exportorv2

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
	"github.com/golang/protobuf/proto"
)

const (
	// 信息所在的行
	TypeSheetRow_Pragma    = 0 // 配置
	TypeSheetRow_Comment   = 1 // 字段名(对应proto)
	TypeSheetRow_DataBegin = 2 // 数据开始
)

const (
	// 信息所在列
	TypeSheetCol_ObjectType = 0 // 对象类型
	TypeSheetCol_FieldName  = 1 // 字段名
	TypeSheetCol_FieldType  = 2 // 字段类型
	TypeSheetCol_Value      = 3 // 值
	TypeSheetCol_Comment    = 4 // 注释
	TypeSheetCol_Meta       = 5 // 特性
)

type TypeSheet struct {
	*Sheet
}

func (self *TypeSheet) Parse(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	// 是否继续读行
	var readingLine bool = true

	var td *model.Descriptor

	rawPragma := self.GetCellData(TypeSheetRow_Pragma, 0)

	if err := proto.UnmarshalText(rawPragma, &localFD.Pragma); err != nil {
		self.Row = TypeSheetRow_Pragma
		self.Column = 0
		log.Errorf("parse pragma failed: %s", rawPragma)
		goto ErrorStop
	}

	if localFD.Pragma.TableName == "" {
		self.Row = TypeSheetRow_Pragma
		self.Column = 0
		log.Errorf("@Types TableName is empty")
		goto ErrorStop
	}

	if localFD.Pragma.Package == "" {
		self.Row = TypeSheetRow_Pragma
		self.Column = 0
		log.Errorf("@Types Package is empty")
		goto ErrorStop
	}

	// 遍历每一行
	for self.Row = TypeSheetRow_DataBegin; readingLine; self.Row++ {

		// ====================解析对象类型====================
		// 第一列是空的，结束
		if self.GetCellData(self.Row, TypeSheetCol_ObjectType) == "" {
			break
		}

		var fd model.FieldDescriptor

		rawTypeName := self.GetCellData(self.Row, TypeSheetCol_ObjectType)

		existType, ok := localFD.DescriptorByName[rawTypeName]

		if ok {

			td = existType

		} else {

			td = model.NewDescriptor()
			td.Name = rawTypeName
			localFD.Add(td)
		}

		// ====================解析字段名====================
		fd.Name = self.GetCellData(self.Row, TypeSheetCol_FieldName)

		// ====================解析字段类型====================
		rawFieldType := self.GetCellData(self.Row, TypeSheetCol_FieldType)

		// 开始在本地symbol中找
		testFD := localFD

		for {

			result := parseFieldType(testFD, rawFieldType, &fd)

			// 找不到就换全局范围找
			if result == parseFieldTypeResult_UnknownFieldType {

				if testFD == localFD {
					testFD = globalFD
					continue
				}

				log.Errorln("unknown field type: ", rawFieldType)

			} else if result == parseFieldTypeResult_OK {
				break
			}

			self.Column = TypeSheetCol_FieldType
			goto ErrorStop

		}

		// ====================解析值====================
		rawValue := self.GetCellData(self.Row, TypeSheetCol_Value)

		var kind model.DescriptorKind

		// 非空值是枚举
		if rawValue != "" {

			// 解析枚举值
			if v, err := strconv.Atoi(rawValue); err == nil {
				fd.EnumValue = int32(v)
			} else {
				self.Column = TypeSheetCol_Value
				log.Errorln("parse type value failed:", err)
				goto ErrorStop
			}
			kind = model.DescriptorKind_Enum
		} else {
			kind = model.DescriptorKind_Struct
		}

		if td.Kind == model.DescriptorKind_None {
			td.Kind = kind
			// 一些字段有填值, 一些没填值
		} else if td.Kind != kind {
			self.Column = TypeSheetCol_Value
			log.Errorln("buildin kind shold be same", td.Kind, kind)
			goto ErrorStop
		}
		// ====================解析注释====================
		fd.Comment = self.GetCellData(self.Row, TypeSheetCol_Comment)

		// ====================解析特性====================
		metaString := self.GetCellData(self.Row, TypeSheetCol_Meta)

		if err := proto.UnmarshalText(metaString, &fd.Meta); err != nil {
			log.Errorln("parse field header failed", err)
			return false
		}

		td.Add(&fd)

	}

	return self.checkProtobufCompatibility(localFD)

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

const (
	parseFieldTypeResult_OK = iota
	parseFieldTypeResult_OtherError
	parseFieldTypeResult_UnknownFieldType
)

func parseFieldType(localFD *model.FileDescriptor, rawFieldType string, fd *model.FieldDescriptor) int {

	// 解析普通类型
	if ft, ok := model.ParseFieldType(rawFieldType); ok {
		fd.Type = ft
	} else {

		// 解析内建类型
		if desc, ok := localFD.DescriptorByName[rawFieldType]; ok {

			// 只有枚举( 结构体不允许再次嵌套, 增加理解复杂度 )
			if desc.Kind != model.DescriptorKind_Enum {
				log.Errorln("struct field can only be normal type and enum", rawFieldType)
				return parseFieldTypeResult_OtherError
			}

			fd.Type = model.FieldType_Enum
			fd.Complex = desc

		} else {

			return parseFieldTypeResult_UnknownFieldType
		}

	}

	return parseFieldTypeResult_OK

}

// 检查protobuf兼容性
func (self *TypeSheet) checkProtobufCompatibility(fileD *model.FileDescriptor) bool {

	for _, bt := range fileD.Descriptors {
		if bt.Kind == model.DescriptorKind_Enum {

			// proto3 需要枚举有0值
			if _, ok := bt.FieldByNumber[0]; !ok {
				log.Errorf("proto3 require enum has value 0 in '%s'", bt.Name)
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
