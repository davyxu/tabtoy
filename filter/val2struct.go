package filter

import (
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
	"github.com/davyxu/tabtoy/proto/tool"
)

// 自定义的token id
const (
	Token_Unknown = iota
	Token_Numeral
	Token_String
	Token_WhiteSpace
	Token_LineEnd
	Token_UnixStyleComment
	Token_Identifier
	Token_True
	Token_False
	Token_Comma
)

func FieldByNameWithMeta(msgD *pbmeta.Descriptor, name string) *pbmeta.FieldDescriptor {

	for i := 0; i < msgD.FieldCount(); i++ {
		fd := msgD.Field(i)

		if fd.Name() == name {
			return fd
		}

		meta := data.GetFieldMeta(fd)

		if meta != nil && meta.Alias == name {

			return fd
		}

	}

	return nil

}

func Value2Struct(meta *tool.FieldMeta, structValue string, fd *pbmeta.FieldDescriptor, callback func(string, string) bool) (isValue2Struct bool, hasError bool) {

	if meta == nil {
		return
	}

	if meta.String2Struct == false {
		return
	}

	if !fd.IsMessageType() {
		hasError = true
		log.Errorf("%s is not message type", fd.Name())
		return
	}

	msgD := fd.MessageDesc()

	if msgD == nil {
		hasError = true
		log.Errorf("%s message not found", fd.Name())
		return
	}

	p := newLineParser(fd, structValue)

	// 匹配顺序从高到低

	defer func() {

		err := recover()

		switch err.(type) {
		// 运行时错误
		case interface {
			RuntimeError()
		}:
			// 继续外抛， 方便调试
			panic(err)
		case error:
			hasError = true
			log.Errorf("field: %s parse error, %v", fd.Name(), err)

		default:
			isValue2Struct = true
		}

	}()

	p.NextToken()

	for {

		if p.TokenID() != Token_Identifier {
			hasError = true
			log.Errorf("expect key in field: %s", fd.Name())
			return
		}

		key := p.TokenValue()

		structFD := FieldByNameWithMeta(msgD, key)

		// 尝试查找字段定义
		if structFD == nil {

			hasError = true
			log.Errorf("%s field not found ", key)
			return
		}

		p.NextToken()

		if p.TokenID() != Token_Comma {
			hasError = true
			log.Errorf("%s need ':' split value", key)
			return
		}

		p.NextToken()

		value := p.TokenValue()

		// 按照正常流程转换值
		if afterValue, ok := ValueConvetor(structFD, value); ok {

			if !callback(key, afterValue) {
				hasError = true
				return
			}

		} else {
			hasError = true
			log.Errorf("%s convert failed", key)
			return
		}

		p.NextToken()

	}

	return
}

func isNumeral(t pbprotos.FieldDescriptorProto_Type) bool {

	switch t {
	case pbprotos.FieldDescriptorProto_TYPE_DOUBLE,
		pbprotos.FieldDescriptorProto_TYPE_FLOAT,
		pbprotos.FieldDescriptorProto_TYPE_INT64,
		pbprotos.FieldDescriptorProto_TYPE_UINT64,
		pbprotos.FieldDescriptorProto_TYPE_INT32,
		pbprotos.FieldDescriptorProto_TYPE_UINT32,
		pbprotos.FieldDescriptorProto_TYPE_FIXED64,
		pbprotos.FieldDescriptorProto_TYPE_FIXED32,
		pbprotos.FieldDescriptorProto_TYPE_SFIXED32,
		pbprotos.FieldDescriptorProto_TYPE_SFIXED64,
		pbprotos.FieldDescriptorProto_TYPE_SINT32,
		pbprotos.FieldDescriptorProto_TYPE_SINT64:
		return true

	}

	return false

}
