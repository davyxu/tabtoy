package exportorv2

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type typeCell struct {
	value string
	col   int
}

// 类型表的数据
type typeModel struct {
	colData map[string]*typeCell

	done bool

	row int

	fd *model.FieldDescriptor

	rawFieldType string
}

func (self *typeModel) getValue(row string) (string, int) {
	if v, ok := self.colData[row]; ok {
		return v.value, v.col
	}

	return "", -1
}

func newTypeModel() *typeModel {
	return &typeModel{
		colData: make(map[string]*typeCell),
		fd:      model.NewFieldDescriptor(),
	}
}

type typeModelRoot struct {
	pragma string

	models []*typeModel

	unknownModel []*typeModel
	fieldTypeCol int

	Col int
	Row int
}

func (self *typeModelRoot) ParsePragma(localFD *model.FileDescriptor) bool {

	if err := localFD.Pragma.Parse(self.pragma); err != nil {
		log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_PragmaParseFailed), self.pragma)
		return false
	}

	if localFD.Pragma.GetString("TableName") == "" {
		log.Errorf("%s", i18n.String(i18n.TypeSheet_TableNameIsEmpty))
		return false
	}

	if localFD.Pragma.GetString("Package") == "" {
		log.Errorf("%s", i18n.String(i18n.TypeSheet_PackageIsEmpty))
		return false
	}

	return true
}

// 解析类型表里的所有类型描述
func (self *typeModelRoot) ParseData(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	var td *model.Descriptor

	reservedRowFieldTypeName := localFD.Pragma.GetString("TableName") + "Define"

	// 每一行
	for _, m := range self.models {

		self.Row = m.row

		var rawTypeName string

		rawTypeName, self.Col = m.getValue("ObjectType")

		if rawTypeName == reservedRowFieldTypeName {
			log.Errorf("%s '%s'", i18n.String(i18n.DataHeader_UseReservedTypeName), rawTypeName)
			return false
		}

		existType, ok := localFD.DescriptorByName[rawTypeName]

		if ok {

			td = existType

		} else {

			td = model.NewDescriptor()
			td.Name = rawTypeName
			localFD.Add(td)
		}

		// 字段名
		m.fd.Name, self.Col = m.getValue("FieldName")

		// 解析类型
		m.rawFieldType, self.Col = m.getValue("FieldType")
		self.fieldTypeCol = self.Col

		fieldType, complexType, ok := findFieldType(localFD, globalFD, m.rawFieldType)
		if !ok {
			return false
		}

		if fieldType == model.FieldType_None {
			self.unknownModel = append(self.unknownModel, m)
		}

		m.fd.Type = fieldType
		m.fd.Complex = complexType

		var rawFieldValue string
		// 解析值
		rawFieldValue, self.Col = m.getValue("Value")

		kind, enumvalue, ok := parseFieldValue(rawFieldValue)
		if !ok {
			return false
		}

		if td.Kind == model.DescriptorKind_None {
			td.Kind = kind
			// 一些字段有填值, 一些没填值
		} else if td.Kind != kind {
			log.Errorf("%s", i18n.String(i18n.TypeSheet_DescriptorKindNotSame))
			return false
		}

		if td.Kind == model.DescriptorKind_Enum {
			if _, ok := td.FieldByNumber[enumvalue]; ok {
				log.Errorf("%s %d", i18n.String(i18n.TypeSheet_DuplicatedEnumValue), enumvalue)
				return false
			}
		}

		m.fd.EnumValue = enumvalue

		m.fd.Comment, self.Col = m.getValue("Comment")

		var rawMeta string
		rawMeta, self.Col = m.getValue("Meta")

		if err := m.fd.Meta.Parse(rawMeta); err != nil {
			log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_FieldMetaParseFailed), err.Error())
			return false
		}

		// 别名
		var rawAlias string
		rawAlias, self.Col = m.getValue("Alias")
		if self.Col != -1 {
			m.fd.Meta.SetString("Alias", rawAlias)
		}

		// 默认值
		var rawDefault string
		rawDefault, self.Col = m.getValue("Default")
		if self.Col != -1 {
			m.fd.Meta.SetString("Default", rawDefault)
		}

		td.Add(m.fd)

	}

	return true
}

func (self *typeModelRoot) SolveUnknownModel(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	for _, m := range self.unknownModel {

		self.Row = m.row
		self.Col = self.fieldTypeCol

		fieldType, complexType, ok := findFieldType(localFD, globalFD, m.rawFieldType)
		if !ok {
			return false
		}

		// 实在是找不到了, 没辙了
		if fieldType == model.FieldType_None {
			log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_FieldTypeNotFound), m.rawFieldType)
			return false
		}

		m.fd.Type = fieldType
		m.fd.Complex = complexType
	}

	return true
}

func findFieldType(localFD *model.FileDescriptor, globalFD *model.FileDescriptor, rawFieldType string) (model.FieldType, *model.Descriptor, bool) {

	// 开始在本地symbol中找
	testFD := localFD

	for {

		fieldType, complexType, ok := findlocalFieldType(testFD, rawFieldType)

		if !ok {
			return model.FieldType_None, nil, false
		}

		if fieldType != model.FieldType_None {
			return fieldType, complexType, true
		}

		// 已经是全局范围, 说明找不到
		if testFD == globalFD {

			return model.FieldType_None, nil, true
		}

		// 找不到就换全局范围找
		testFD = globalFD
	}

}

// bool表示是否有错
func findlocalFieldType(localFD *model.FileDescriptor, rawFieldType string) (model.FieldType, *model.Descriptor, bool) {

	// 解析普通类型
	if ft, ok := model.ParseFieldType(rawFieldType); ok {

		return ft, nil, true

	}

	// 解析内建类型
	if desc, ok := localFD.DescriptorByName[rawFieldType]; ok {

		// 只有枚举( 结构体不允许再次嵌套, 增加理解复杂度 )
		if desc.Kind != model.DescriptorKind_Enum {
			log.Errorf("%s, '%s'", i18n.String(i18n.TypeSheet_StructFieldCanNotBeStruct), rawFieldType)

			return model.FieldType_None, nil, false
		}

		return model.FieldType_Enum, desc, true

	}

	// 没找到类型, 待二次pass
	return model.FieldType_None, nil, true

}

func parseFieldValue(rawFieldValue string) (model.DescriptorKind, int32, bool) {

	// 非空值是枚举
	if rawFieldValue != "" {

		v, err := strconv.Atoi(rawFieldValue)
		// 解析枚举值
		if err != nil {

			log.Errorf("%s, %s", i18n.String(i18n.TypeSheet_EnumValueParseFailed), err.Error())
			return model.DescriptorKind_None, 0, false
		}

		return model.DescriptorKind_Enum, int32(v), true
	}

	return model.DescriptorKind_Struct, 0, true
}
