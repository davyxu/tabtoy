package printer

import (
	"strconv"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

func valueWrapperJson(t model.FieldType, node *model.Node) string {

	switch t {
	case model.FieldType_String:
		return util.StringEscape(node.Value)
	case model.FieldType_Enum:
		return strconv.Itoa(int(node.EnumValue))
	}

	return node.Value
}

func PrintJson(tab *model.Table, rootName string, version string, outfile string) bool {

	var fp FilePrinter

	fp.Printf("{\n")

	fp.Printf("	\"Tool\": \"github.com/davyxu/tabtoy\",\n")
	fp.Printf("	\"Version\": \"%s\",\n", version)

	fp.Printf("	\"%s\":[\n", rootName)

	// 遍历每一行
	for rIndex, r := range tab.Recs {

		fp.Printf("		{ ")

		// 遍历每一列
		for rootFieldIndex, node := range r.Nodes {

			if node.IsRepeated {
				fp.Printf("\"%s\":[ ", node.Name)
			} else {
				fp.Printf("\"%s\": ", node.Name)
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				if node.IsRepeated {

					// repeated 值序列
					for arrIndex, valueNode := range node.Child {

						fp.Printf("%s", valueWrapperJson(node.Type, valueNode))

						// 多个值分割
						if arrIndex < len(node.Child)-1 {
							fp.Printf(", ")
						}

					}
				} else {
					// 单值
					valueNode := node.Child[0]

					fp.Printf("%s", valueWrapperJson(node.Type, valueNode))

				}

			} else {

				// 遍历repeated的结构体
				for structIndex, structNode := range node.Child {

					// 结构体开始
					fp.Printf("{ ")

					// 遍历一个结构体的字段
					for structFieldIndex, fieldNode := range structNode.Child {

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						fp.Printf("\"%s\": %s", fieldNode.Name, valueWrapperJson(fieldNode.Type, valueNode))

						// 结构体字段分割
						if structFieldIndex < len(structNode.Child)-1 {
							fp.Printf(", ")
						}

					}

					// 结构体结束
					fp.Printf(" }")

					// 多个结构体分割
					if structIndex < len(node.Child)-1 {
						fp.Printf(", ")
					}

				}

			}

			if node.IsRepeated {
				fp.Printf(" ]")
			}

			// 根字段分割
			if rootFieldIndex < len(r.Nodes)-1 {
				fp.Printf(", ")
			}

		}

		fp.Printf(" }")

		if rIndex < len(tab.Recs)-1 {
			fp.Printf(",")
		}

		fp.Printf("\n")

	}

	fp.Printf("	]\n")
	fp.Printf("}")

	return fp.Write(outfile)

}
