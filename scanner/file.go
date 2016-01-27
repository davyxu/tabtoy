package scanner

import (
	"path"
	"strings"

	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/golang/protobuf/proto"
	"github.com/tealeg/xlsx"
)

type File struct {
	Raw      *xlsx.File
	Sheets   []*Sheet // 多个sheet
	SheetMap map[string]*Sheet
	descpool *pbmeta.DescriptorPool // 协议描述池
	Name     string                 // 电子表格的名称(文件名无后缀)
	FileName string
}

func (self *File) Open(filename string) bool {

	self.Name = makeTableName(filename)
	self.FileName = filename

	var err error
	self.Raw, err = xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return false
	}

	// 这里将所有sheet表都合并导出到一个pbt
	for _, sheet := range self.Raw.Sheets {

		self.Add(sheet)
	}

	return true
}

func (self *File) Add(sheet *xlsx.Sheet) *Sheet {
	header := getHeader(sheet)

	// 没有找到导出头,忽略
	if header == nil {
		return nil
	}

	mySheet := newSheet(self, sheet, header)

	// TODO 添加命令行导出忽略
	self.Sheets = append(self.Sheets, mySheet)

	self.SheetMap[sheet.Name] = mySheet

	return mySheet
}

// 制作表名
func makeTableName(filename string) string {
	baseName := path.Base(filename)
	return strings.TrimSuffix(baseName, path.Ext(baseName))
}

func NewFile(descpool *pbmeta.DescriptorPool) *File {

	self := &File{
		Sheets:   make([]*Sheet, 0),
		SheetMap: make(map[string]*Sheet),
		descpool: descpool,
		Raw:      xlsx.NewFile(),
	}

	return self
}

func getHeader(sheet *xlsx.Sheet) *tool.ExportHeader {

	headerString := strings.TrimSpace(sheet.Cell(0, 0).Value)

	if headerString == "" {
		return nil
	}

	var header tool.ExportHeader

	if err := proto.UnmarshalText(headerString, &header); err != nil {
		return nil
	}

	return &header
}
