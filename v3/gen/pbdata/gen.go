package pbdata

import (
	"fmt"
	"github.com/davyxu/tabtoy/v3/model"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
	"io/ioutil"
)

func exportTable(globals *model.Globals, pbFile protoreflect.FileDescriptor, tab *model.DataTable, combineRoot *dynamicpb.Message) {
	md := pbFile.Messages().ByName(protoreflect.Name(tab.OriginalHeaderType))

	combineField := combineRoot.Descriptor().Fields().ByName(protoreflect.Name(tab.OriginalHeaderType))
	list := combineRoot.NewField(combineField).List()

	// 每个表的所有列
	headers := globals.Types.AllFieldByName(tab.OriginalHeaderType)

	// 遍历每一行
	for row := 1; row < len(tab.Rows); row++ {

		rowMsg := dynamicpb.NewMessage(md)

		for col, field := range headers {

			if globals.CanDoAction(model.ActionNoGenPbBinary, field) {
				continue
			}

			fd := md.Fields().ByName(protoreflect.Name(field.FieldName))

			// 在单元格找到值
			valueCell := tab.GetCell(row, col)
			if valueCell == nil {
				continue
			}

			if field.IsArray() {
				list := rowMsg.NewField(fd).List()
				tableValue2PbValueList(globals, valueCell, field, list)
				rowMsg.Set(fd, protoreflect.ValueOfList(list))
			} else {
				pbValue := tableValue2PbValue(globals, valueCell.Value, field)
				rowMsg.Set(fd, pbValue)
			}

		}

		list.Append(protoreflect.ValueOf(rowMsg))
	}

	combineRoot.Set(combineField, protoreflect.ValueOfList(list))
}

func Generate(globals *model.Globals) (data []byte, err error) {

	pbFile, err := buildDynamicType(globals)
	if err != nil {
		return nil, err
	}

	combineD := pbFile.Messages().ByName(protoreflect.Name(globals.CombineStructName))

	combineRoot := dynamicpb.NewMessage(combineD)

	// 所有的表
	for _, tab := range globals.Datas.AllTables() {
		exportTable(globals, pbFile, tab, combineRoot)
	}

	return proto.Marshal(combineRoot)
}

func Output(globals *model.Globals, param string) (err error) {

	pbFile, err := buildDynamicType(globals)
	if err != nil {
		return err
	}

	for _, tab := range globals.Datas.AllTables() {

		combineD := pbFile.Messages().ByName(protoreflect.Name(globals.CombineStructName))

		combineRoot := dynamicpb.NewMessage(combineD)

		exportTable(globals, pbFile, tab, combineRoot)

		data, err := proto.Marshal(combineRoot)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/%s.pbb", param, tab.HeaderType), data, 0666)

		if err != nil {
			return err
		}
	}

	return nil
}
