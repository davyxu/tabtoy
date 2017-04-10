package printer

import (
	"encoding/json"

	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type typePrinter struct {
}

// 一个列字段
type typeFieldModel struct {
	Name       string
	Type       string
	Kind       string
	IsRepeated bool
	Meta       map[string]interface{}
	Comment    string

	Value int
}

// 一张表的类型信息
type typeStructModel struct {
	Name   string
	Fields []*typeFieldModel
}

// 整个文件类型信息
type typeFileModel struct {
	Tool    string
	Version string
	Structs []*typeStructModel
	Enums   []*typeStructModel
}

func (self *typePrinter) Run(g *Globals) *Stream {

	bf := NewStream()

	var fm typeFileModel
	fm.Tool = "github.com/davyxu/tabtoy"
	fm.Version = g.Version

	// 遍历所有类型
	for _, d := range g.FileDescriptor.Descriptors {

		// 这给被限制输出
		if !d.File.MatchTag(".type") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), d.Name)
			continue
		}

		structM := &typeStructModel{
			Name: d.Name,
		}

		// 遍历字段
		for _, fd := range d.Fields {

			// 对CombineStruct的XXDefine对应的字段
			if d.Usage == model.DescriptorUsage_CombineStruct {

				// 这个字段被限制输出
				if fd.Complex != nil && !fd.Complex.File.MatchTag(".type") {
					continue
				}
			}

			field := &typeFieldModel{
				Name:       fd.Name,
				Type:       fd.TypeString(),
				Kind:       fd.KindString(),
				IsRepeated: fd.IsRepeated,
				Comment:    fd.Comment,
				Meta:       fd.Meta.Raw(),
			}

			switch d.Kind {
			case model.DescriptorKind_Struct:
				field.Value = 0
			case model.DescriptorKind_Enum:
				field.Value = int(fd.EnumValue)
			}

			structM.Fields = append(structM.Fields, field)

		}

		switch d.Kind {
		case model.DescriptorKind_Struct:
			fm.Structs = append(fm.Structs, structM)
		case model.DescriptorKind_Enum:
			fm.Enums = append(fm.Enums, structM)
		}

	}

	data, err := json.MarshalIndent(&fm, "", " ")
	if err != nil {
		log.Errorln(err)
	}

	bf.WriteBytes(data)

	return bf
}

func init() {

	RegisterPrinter("type", &typePrinter{})

}
