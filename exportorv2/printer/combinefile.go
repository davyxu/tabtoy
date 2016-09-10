package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type CombineFile struct {
	*BinaryFile                         // 最终输出文件
	nameChecker map[string]*model.Table //  防止table重名
	tableCount  int32                   // 添加了多少表

	rootTypes   []*model.FieldDefine // 根类型, XXDefine
	packageName string               // 最终的包

	indexes      []tableIndex
	indexChecker map[tableIndex]bool
}

type tableIndexType int

const (
	tableIndexType_None tableIndexType = iota
	tableIndexType_Lua
	tableIndexType_Sharp
)

type tableIndex struct {
	Raw   *model.FieldDefine
	Index *model.FieldDefine
	Type  tableIndexType
}

// 合并每个表带的类型
func (self *CombineFile) CombineType(inputFile string, tts *model.BuildInTypeSet) bool {
	// 有表格里描述的包名不一致, 无法合成最终的文件
	if self.packageName != "" && self.packageName != tts.Pragma.Package {
		log.Errorf("combine file 'Package' in @Types diff: %s", inputFile)
		return false
	}

	self.packageName = tts.Pragma.Package

	fileType := tts.FileType.Fields[0]
	// 收集根文件类型, 以做最终合并时生成
	self.rootTypes = append(self.rootTypes, fileType)

	for _, bt := range tts.Types {

		// 非行类型的, 全部忽略
		if bt.Usage != model.BuildIntTypeUsage_RowType {
			continue
		}

		for _, field := range bt.Fields {

			key := tableIndex{
				Raw:   fileType,
				Index: field,
			}

			if field.Meta.CSharpIndex {
				key.Type = tableIndexType_Sharp

				if _, ok := self.indexChecker[key]; ok {
					log.Errorf("duplicate CSharpIndex field: %s", field.String())
					return false
				}

				self.indexChecker[key] = true
				self.indexes = append(self.indexes, key)
			}

		}

	}

	return true
}

// 输出带有合并类型的C#文件
func (self *CombineFile) PrintCombineCSharp(toolVersion string, name string) *BinaryFile {

	combineFileTypeSet := model.NewBuildInTypeSet()
	combineFileTypeSet.Pragma.Package = self.packageName
	fileStruct := model.NewBuildInType()
	fileStruct.Name = name
	fileStruct.Kind = model.BuildInTypeKind_Struct

	combineFileTypeSet.Add(fileStruct)

	// 遍历所有导出root类型, 添加到XXFile的字段上
	for index, ft := range self.rootTypes {

		ft.Order = int32(index)

		fileStruct.Add(ft)
	}

	return PrintCSharp(combineFileTypeSet, self.indexes, toolVersion)
}

func (self *CombineFile) WriteBinary(tab *model.Table) bool {

	if _, ok := self.nameChecker[tab.Name]; ok {
		log.Errorln("duplicate table name in combine binary output:", tab.Name)
		return false
	}

	self.nameChecker[tab.Name] = tab

	// 表所在的字段
	self.WriteInt32(model.MakeTag(model.FieldType_Struct, self.tableCount))
	self.tableCount++

	self.WriteInt32(int32(len(tab.Recs)))

	// 遍历每一行
	for _, r := range tab.Recs {

		// 遍历每一列
		for _, node := range r.Nodes {

			// 写入字段索引
			self.WriteInt32(node.Tag())

			// 写入数量
			if node.IsRepeated {
				self.WriteInt32(int32(len(node.Child)))
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				for _, valueNode := range node.Child {

					self.WriteNodeValue(node.Type, valueNode)
				}

			} else {

				// 遍历repeated的结构体
				for _, structNode := range node.Child {

					// 遍历一个结构体的字段
					for _, fieldNode := range structNode.Child {

						// 写入字段索引
						self.WriteInt32(fieldNode.Tag())

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						self.WriteNodeValue(fieldNode.Type, valueNode)

					}

				}

			}

		}

	}

	return true

}

const combineFileVersion = 1

func NewCombineFile() *CombineFile {
	self := &CombineFile{
		nameChecker:  make(map[string]*model.Table),
		indexChecker: make(map[tableIndex]bool),
		BinaryFile:   NewBinaryFile("Combine"),
	}

	self.WriteString("TABTOY")
	self.WriteInt32(combineFileVersion)

	return self
}
