package exportor

import (
	"strings"

	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/exportorv1/data"
	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/davyxu/tabtoy/util"
	"github.com/tealeg/xlsx"
)

const (
	// 信息所在的行
	DataIndex_ExportHeader = 0 // 导出头
	DataIndex_FieldName    = 1 // 字段名(对应proto)
	DataIndex_FieldDesc    = 2 // 字段描述
	DataIndex_DataBegin    = 3 // 数据开始
)

type Sheet struct {
	*xlsx.Sheet
	header *tool.ExportHeaderV1

	cursor int // 当前行

	index int // 当前列

	file *File // 指向父级

	FieldHeader []string // 有效的字段行，可以做多sheet对比
}

// 获取单元格
func (self *Sheet) GetCellData(cursor, index int) string {

	return strings.TrimSpace(self.Cell(cursor, index).Value)
}

// 设置单元格
func (self *Sheet) SetCellData(cursor, index int, data string) {

	self.Cell(cursor, index).Value = data
}

// 检查字段行的长度
func (self *Sheet) ParseProtoField() bool {

	// proto字段导引头

	for index := 0; ; index++ {
		fieldName := self.GetCellData(DataIndex_FieldName, index)

		if fieldName == "" {
			break
		}

		self.FieldHeader = append(self.FieldHeader, fieldName)
	}

	// 没有导引头
	return len(self.FieldHeader) > 0
}

func (self *Sheet) checkProtoHeader() (*data.DynamicMessage, *pbmeta.Descriptor, *pbmeta.FieldDescriptor) {

	// 指定的导出文件类型获得描述
	fileDesc := self.file.descpool.MessageByFullName(self.header.ProtoTypeName)

	if fileDesc == nil {
		log.Errorf("can not found record descriptor, '%s'", self.header.ProtoTypeName)
		return nil, nil, nil
	}

	lineMsgDesc := fileDesc.FieldByName(self.header.RowFieldName)

	// 找不到定义
	if lineMsgDesc == nil {
		log.Errorf("row field type not found: %s", self.header.RowFieldName)
		return nil, nil, nil
	}

	// 行描述类型必须是数组
	if !lineMsgDesc.IsRepeated() {
		log.Errorf("row field type must be repeated type: %s", self.header.RowFieldName)
		return nil, nil, nil
	}

	// 根据描述创建输出文件消息及每一行的消息结构类型
	return data.NewDynamicMessage(fileDesc), lineMsgDesc.MessageDesc(), lineMsgDesc
}

type RecordInfo struct {
	FieldName string // 表格中的字段描述,可能是a.b.c
	CellValue string // 表格中的值
	FieldDesc *pbmeta.FieldDescriptor
	FieldMsg  *data.DynamicMessage // 这个字段所在的Message
	FieldMeta *tool.FieldMetaV1    // 扩展字段
}

func (self *RecordInfo) Value() string {

	if self.CellValue == "" && self.FieldMeta != nil {
		return self.FieldMeta.DefaultValue
	}

	return self.CellValue
}

func (self *Sheet) IterateData(callback func(*RecordInfo) bool) (*data.DynamicMessage, bool) {

	// 检查引导头
	if !self.ParseProtoField() {
		return nil, true
	}

	// 是否继续读行
	var readingLine bool = true

	// 检查引导Proto字段
	sheetMsg, rowMsgDesc, lineFieldDesc := self.checkProtoHeader()

	if sheetMsg == nil {
		goto ErrorStop
	}

	// 遍历每一行
	for self.cursor = DataIndex_DataBegin; readingLine; self.cursor++ {

		// 第一列是空的，结束
		if self.GetCellData(self.cursor, 0) == "" {
			break
		}

		lineMsg := data.NewDynamicMessage(rowMsgDesc)

		if lineMsg == nil {
			break
		}

		// 遍历每一列
		for self.index = 0; self.index < len(self.FieldHeader); self.index++ {

			ri := new(RecordInfo)

			// Proto字段头
			ri.FieldName = self.FieldHeader[self.index]

			// 原始值
			ri.CellValue = self.GetCellData(self.cursor, self.index)

			// #开头表示注释, 跳过
			if strings.Index(ri.FieldName, "#") == 0 {
				continue
			}

			ri.FieldMsg, ri.FieldDesc = makeCompactAccessor(ri.FieldName, lineMsg)

			// 字段匹配错误
			if ri.FieldMsg == nil || ri.FieldDesc == nil {
				goto ErrorStop
			}

			var ok bool
			// 取扩展元信息
			ri.FieldMeta, ok = data.GetFieldMeta(ri.FieldDesc)

			if !ok {
				goto ErrorStop
			}

			if data.DebuggingLevel >= 1 {
				r, c := self.GetRC()
				log.Debugf("(%s) %s=%s", util.ConvR1C1toA1(r, c), ri.FieldName, ri.Value)
			}

			if !callback(ri) {
				goto ErrorStop
			}

		}

		sheetMsg.AddRepeatedMessage(lineFieldDesc, lineMsg)

	}

	return sheetMsg, true

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return nil, false
}

// 取行列信息
func (self *Sheet) GetRC() (int, int) {

	return self.cursor + 1, self.index + 1

}

// 解析a.b.c字段，由给定的msg找到这些字段病返回字段访问器
func makeCompactAccessor(compactFieldName string, inputMsg *data.DynamicMessage) (*data.DynamicMessage, *pbmeta.FieldDescriptor) {

	// 将路径按点分割
	fieldNameList := strings.Split(compactFieldName, ".")

	msg := inputMsg

	for _, fieldName := range fieldNameList {

		// 这个路径对应的描述器
		fd := msg.Desc.FieldByName(fieldName)
		if fd == nil {
			log.Errorf("field type name not found, %s in %s", fieldName, compactFieldName)
			return nil, nil
		}

		// 消息进行添加并迭代
		if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {

			// 字段中带repeated message的不支持, 使用string2struct解析字段
			fdmeta, ok := data.GetFieldMeta(fd)
			if !ok {
				return nil, nil
			}

			if fd.IsRepeated() && fdmeta != nil && !fdmeta.String2Struct {
				log.Errorf("DO NOT support repeated message field, use 'string2struct' instead")
				return nil, nil
			}

			existMsg := msg.GetMessage(fd)

			if existMsg == nil {

				newMsg := data.NewDynamicMessage(fd.MessageDesc())

				// 重复消息字段
				if fd.IsRepeated() {

					if !msg.AddRepeatedMessage(fd, newMsg) {
						return nil, nil
					}

				} else { //普通值字段

					if !msg.SetMessage(fd, newMsg) {
						return nil, nil
					}

				}

				msg = newMsg
			} else {
				msg = existMsg
			}

		} else {
			// 非消息返回当前层的反射器

			return msg, fd
		}

	}

	// 纯消息, 字段使用父级
	return msg, inputMsg.Desc.FieldByName(compactFieldName)
}

func newSheet(file *File, sheet *xlsx.Sheet, header *tool.ExportHeaderV1) *Sheet {
	self := &Sheet{
		file:        file,
		Sheet:       sheet,
		header:      header,
		FieldHeader: make([]string, 0),
	}

	return self
}
