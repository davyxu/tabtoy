package printer

import (
	"sync"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

type TableIndex struct {
	Index *model.FieldDescriptor // 表头里的索引
	Row   *model.FieldDescriptor // 索引的数据
}

type Globals struct {
	Version       string
	InputFileList []interface{}
	ParaMode      bool
	ProtoVersion  int

	Printers []*PrinterContext

	CombineStructName string // 不包含路径, 用作

	*model.FileDescriptor // 类型信息.用于添加各种导出结构

	tableByName map[string]*model.Table //  防止table重名
	Tables      []*model.Table          // 数据信息.表格数据

	GlobalIndexes []TableIndex      // 类型信息.全局索引
	CombineStruct *model.Descriptor // 类型信息.Combine结构体

	guard sync.Mutex
}

func (self *Globals) PreExport() bool {

	// 当合并结构名没有指定时, 对于代码相关的输出器, 要报错
	if self.CombineStructName == "" && self.hasAnyPrinter(".proto", ".cs") {
		log.Errorf("please specify 'combinename' params while code generating")
		return false
	}

	// 添加XXConfig全局结构
	self.CombineStruct.Name = self.CombineStructName
	self.CombineStruct.Kind = model.DescriptorKind_Struct
	self.CombineStruct.Usage = model.DescriptorUsage_CombineStruct
	self.FileDescriptor.Name = self.CombineStructName
	self.FileDescriptor.Add(self.CombineStruct)
	return true
}

func (self *Globals) hasAnyPrinter(exts ...string) bool {
	for _, ext := range exts {

		for _, p := range self.Printers {
			if p.ext == ext {
				return true
			}
		}
	}

	return false
}

func (self *Globals) AddOutputType(ext string, outfile string) {

	if p, ok := printerByExt[ext]; ok {
		self.Printers = append(self.Printers, &PrinterContext{
			p:       p,
			outFile: outfile,
			ext:     ext,
		})
	} else {
		panic("output type not found:" + ext)
	}

}

func (self *Globals) Print() bool {

	log.Infoln("==========Merge Combined Data==========")

	for _, p := range self.Printers {

		if !p.Start(self) {
			return false
		}
	}

	return true

}

func (self *Globals) AddTypes(localFD *model.FileDescriptor) bool {

	// 将行定义结构也添加到文件中
	for _, d := range localFD.Descriptors {
		self.FileDescriptor.Add(d)
	}

	return true
}

// 合并每个表带的类型
func (self *Globals) AddContent(tab *model.Table) bool {

	localFD := tab.LocalFD

	self.guard.Lock()

	defer self.guard.Unlock()

	// 有表格里描述的包名不一致, 无法合成最终的文件
	if self.Pragma.Package == "" {
		self.Pragma.Package = localFD.Pragma.Package
	} else if self.Pragma.Package != localFD.Pragma.Package {
		log.Errorf("keep all type in same package! %s, %s", self.Pragma.Package, localFD.Pragma.Package)
		return false
	}

	if _, ok := self.tableByName[localFD.Name]; ok {
		log.Errorln("duplicate table name in combine binary output:", localFD.Name)
		return false
	}

	// 表的全局类型信息与合并信息一致
	tab.GlobalFD = self.FileDescriptor

	self.tableByName[localFD.Name] = tab
	self.Tables = append(self.Tables, tab)

	// 每个表在结构体里的字段
	var rowFD model.FieldDescriptor
	rowFD.Name = localFD.Name
	rowFD.Type = model.FieldType_Struct
	rowFD.Complex = localFD.RowDescriptor()
	rowFD.IsRepeated = true
	rowFD.Order = int32(len(self.CombineStruct.Fields) + 1)
	rowFD.Comment = localFD.Name
	self.CombineStruct.Add(&rowFD)

	for _, d := range localFD.Descriptors {

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

func NewGlobals() *Globals {
	self := &Globals{
		tableByName:    make(map[string]*model.Table),
		FileDescriptor: model.NewFileDescriptor(),
		CombineStruct:  model.NewDescriptor(),
	}

	return self
}
