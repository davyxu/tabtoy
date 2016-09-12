package printer

import (
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type TableIndex struct {
	Index *model.FieldDescriptor // 表头里的索引
	Row   *model.FieldDescriptor // 索引的数据
}

type Globals struct {
	Version       string
	InputFileList []string
	ParaMode      bool
	ProtoVersion  int

	Printers []*PrinterContext

	CombineStructName string // 不包含路径, 用作

	*model.FileDescriptor // 类型信息.用于添加各种导出结构

	tableByName map[string]*model.Table //  防止table重名
	Tables      []*model.Table          // 数据信息.表格数据

	GlobalIndexes []TableIndex      // 类型信息.全局索引
	CombineStruct *model.Descriptor // 类型信息.Combine结构体
}

func (self *Globals) AddOutputType(ext string, outdir string) {

	if p, ok := printerByExt[ext]; ok {
		self.Printers = append(self.Printers, &PrinterContext{
			p:      p,
			outDir: outdir,
			ext:    ext,
		})
	} else {
		panic("output type not found:" + ext)
	}

}

func (self *Globals) Run() bool {

	for _, p := range self.Printers {

		if !p.Start(self) {
			return false
		}
	}

	return true

}

// 合并每个表带的类型
func (self *Globals) CombineType(inputFile string, fileD *model.FileDescriptor) bool {
	// 有表格里描述的包名不一致, 无法合成最终的文件

	if self.Pragma.Package == "" {
		self.Pragma.Package = fileD.Pragma.Package
	} else if self.Pragma.Package != fileD.Pragma.Package {
		log.Errorf("keep all type in same package! %s, %s", self.Pragma.Package, fileD.Pragma.Package)
		return false
	}

	// 每个表在结构体里的字段
	var rowFD model.FieldDescriptor
	rowFD.Name = fileD.Name
	rowFD.Type = model.FieldType_Struct
	rowFD.Complex = fileD.RowDescriptor()
	rowFD.IsRepeated = true
	rowFD.Order = int32(len(self.CombineStruct.Fields) + 1)
	rowFD.Comment = fileD.Name
	self.CombineStruct.Add(&rowFD)

	// 将行定义结构也添加到文件中

	for _, d := range fileD.Descriptors {
		self.FileDescriptor.Add(d)
	}

	for _, d := range fileD.Descriptors {

		// 非行类型的, 全部忽略
		if d.Usage != model.DescriptorUsage_RowType {
			continue
		}

		for _, indexFD := range d.Indexes {

			key := TableIndex{
				Row:   &rowFD,
				Index: indexFD,
			}

			self.GlobalIndexes = append(self.GlobalIndexes, key)

		}

	}

	return true
}

func (self *Globals) CombineData(tab *model.Table) bool {

	if _, ok := self.tableByName[tab.Name]; ok {
		log.Errorln("duplicate table name in combine binary output:", tab.Name)
		return false
	}

	self.tableByName[tab.Name] = tab
	self.Tables = append(self.Tables, tab)

	return true
}

func NewGlobals() *Globals {
	self := &Globals{
		tableByName:    make(map[string]*model.Table),
		FileDescriptor: model.NewFileDescriptor(),
		CombineStruct:  model.NewDescriptor(),
	}

	// 添加XXConfig全局结构
	self.CombineStruct.Name = self.CombineStructName
	self.CombineStruct.Kind = model.DescriptorKind_Struct
	self.CombineStruct.Usage = model.DescriptorUsage_CombineStruct
	self.FileDescriptor.Add(self.CombineStruct)

	return self
}
