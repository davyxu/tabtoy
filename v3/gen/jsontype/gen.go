package jsontype

import (
	"encoding/json"
	"github.com/davyxu/tabtoy/v3/model"
	"sort"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	objByType := map[string]*Object{}

	for _, def := range globals.Types.AllFields() {
		obj := objByType[def.ObjectType]
		if obj == nil {
			obj = &Object{}
			obj.Name = def.ObjectType

			for _, indexDef := range globals.IndexList {
				if indexDef.TableType == def.ObjectType {
					obj.Tags = append(obj.Tags, indexDef.Tags...)
				}
			}

			switch def.Kind {
			case model.TypeUsage_HeaderStruct:
				obj.Type = "Struct"
			case model.TypeUsage_Enum:
				obj.Type = "Enum"
			}
			objByType[def.ObjectType] = obj
		}

		var fd Field
		fd.Name = def.FieldName
		fd.Type = def.FieldType
		fd.Comment = def.Name
		fd.Value = def.Value
		fd.MakeIndex = def.MakeIndex
		fd.ArraySplitter = def.ArraySplitter
		fd.Tags = def.Tags

		obj.Fields = append(obj.Fields, &fd)
	}

	var f File
	f.Version = globals.Version
	f.Tool = "github.com/davyxu/tabtoy"

	for _, obj := range objByType {
		f.Objects = append(f.Objects, obj)
	}

	sort.Slice(f.Objects, func(i, j int) bool {
		return f.Objects[i].Compare(f.Objects[j])
	})

	return json.MarshalIndent(&f, "", "\t")
}
