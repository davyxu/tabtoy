package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/model"
)

func PrintCombineCSharp(filetypes []*model.FieldDefine, toolVersion string, name string, namespace string) *BinaryFile {

	combineFileTypeSet := model.NewBuildInTypeSet()
	combineFileTypeSet.Pragma.Package = namespace
	fileStruct := model.NewBuildInType()
	fileStruct.Name = name
	fileStruct.Kind = model.BuildInTypeKind_Struct

	combineFileTypeSet.Add(fileStruct)

	// 遍历所有导出root类型, 添加到XXFile的字段上
	for index, ft := range filetypes {

		ft.Order = int32(index)

		fileStruct.Add(ft)
	}

	return PrintCSharp(combineFileTypeSet, toolVersion)
}
