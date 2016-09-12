package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type CombineFile struct {
	*BinaryFile // 最终输出文件

	fd *model.FileDescriptor

	nameChecker map[string]*model.Table //  防止table重名
	tableCount  int32                   // 添加了多少表

	indexes       []tableIndex
	combineStruct *model.Descriptor
}

type tableIndex struct {
	Index *model.FieldDescriptor // 表头里的索引
	Row   *model.FieldDescriptor // 索引的数据
}

// 合并每个表带的类型
func (self *CombineFile) CombineType(inputFile string, fileD *model.FileDescriptor) bool {
	// 有表格里描述的包名不一致, 无法合成最终的文件

	if self.fd.Pragma.Package == "" {
		self.fd.Pragma.Package = fileD.Pragma.Package
	} else if self.fd.Pragma.Package != fileD.Pragma.Package {
		log.Errorf("keep all type in same package! %s, %s", self.fd.Pragma.Package, fileD.Pragma.Package)
		return false
	}

	// 每个表在结构体里的字段
	var rowFD model.FieldDescriptor
	rowFD.Name = fileD.Name
	rowFD.Type = model.FieldType_Struct
	rowFD.Complex = fileD.RowDescriptor()
	rowFD.IsRepeated = true
	rowFD.Order = int32(len(self.combineStruct.Fields) + 1)
	rowFD.Comment = fileD.Name
	self.combineStruct.Add(&rowFD)

	// 将行定义结构也添加到文件中

	for _, d := range fileD.Descriptors {
		self.fd.Add(d)
	}

	for _, d := range fileD.Descriptors {

		// 非行类型的, 全部忽略
		if d.Usage != model.DescriptorUsage_RowType {
			continue
		}

		for _, indexFD := range d.Indexes {

			key := tableIndex{
				Row:   &rowFD,
				Index: indexFD,
			}

			self.indexes = append(self.indexes, key)

		}

	}

	return true
}

// 输出带有合并类型的C#文件
func (self *CombineFile) PrintCombineCSharp(toolVersion string) *BinaryFile {

	return PrintCSharp(self.fd, self.indexes, toolVersion)
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

func NewCombineFile(stuctname string) *CombineFile {
	self := &CombineFile{
		nameChecker:   make(map[string]*model.Table),
		BinaryFile:    NewBinaryFile("Combine"),
		fd:            model.NewFileDescriptor(),
		combineStruct: model.NewDescriptor(),
	}

	// 添加XXConfig全局结构
	self.combineStruct.Name = stuctname
	self.combineStruct.Kind = model.DescriptorKind_Struct
	self.combineStruct.Usage = model.DescriptorUsage_CombineStruct
	self.fd.Add(self.combineStruct)

	self.WriteString("TABTOY")
	self.WriteInt32(combineFileVersion)

	return self
}
