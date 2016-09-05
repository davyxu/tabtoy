package exportor

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

		if err := self.Add(sheet); err != nil {
			return false
		}
	}

	return true
}

func (self *File) Add(sheet *xlsx.Sheet) error {
	header, err := getHeader(sheet)

	if err != nil {
		log.Errorf("invalid proto header in file %s: %s", self.FileName, err)
		return err
	}

	// 没有找到导出头,忽略
	if header == nil {
		return nil
	}

	mySheet := newSheet(self, sheet, header)

	// TODO 添加命令行导出忽略
	self.Sheets = append(self.Sheets, mySheet)

	self.SheetMap[sheet.Name] = mySheet

	return nil
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

func getHeader(sheet *xlsx.Sheet) (*tool.ExportHeader, error) {

	headerString := strings.TrimSpace(sheet.Cell(0, 0).Value)

	// 可能是空的sheet
	if headerString == "" {
		return nil, nil
	}

	var header tool.ExportHeader

	// 有可能的字符,一定是头
	if strings.Contains(headerString, "ProtoTypeName") ||
		strings.Contains(headerString, "RowFieldName") {
		if err := proto.UnmarshalText(headerString, &header); err != nil {

			return nil, err
		}
	} else {
		// 有字符, 但并不是头
		return nil, nil
	}

	return &header, nil
}
