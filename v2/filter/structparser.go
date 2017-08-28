package filter

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

// 自定义的token id
const (
	Token_EOF = iota
	Token_WhiteSpace
	Token_LineEnd
	Token_UnixStyleComment
	Token_Identifier
	Token_Numeral
	Token_String
	Token_Comma
	Token_Unknown
)

type structParser struct {
	*golexer.Parser
}

func (self *structParser) Run(fd *model.FieldDescriptor, callback func(string, string) bool) (ok bool) {

	defer golexer.ErrorCatcher(func(err error) {

		log.Errorf("%s, '%s' '%v'", i18n.String(i18n.StructParser_LexerError), fd.Name, err.Error())
	})

	self.NextToken()

	for self.TokenID() != Token_EOF {

		if self.TokenID() != Token_Identifier {
			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_ExpectField), fd.Name)
			return false
		}

		key := self.TokenValue()

		self.NextToken()

		if self.TokenID() != Token_Comma {
			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_UnexpectedSpliter), key)
			return false
		}

		self.NextToken()

		value := self.TokenValue()

		if !callback(key, value) {
			return false
		}

		self.NextToken()

	}

	return true
}

func newStructParser(value string) *structParser {
	l := golexer.NewLexer()

	l.AddMatcher(golexer.NewNumeralMatcher(Token_Numeral))
	l.AddMatcher(golexer.NewStringMatcher(Token_String))

	l.AddIgnoreMatcher(golexer.NewWhiteSpaceMatcher(Token_WhiteSpace))
	l.AddIgnoreMatcher(golexer.NewLineEndMatcher(Token_LineEnd))
	l.AddIgnoreMatcher(golexer.NewUnixStyleCommentMatcher(Token_UnixStyleComment))

	l.AddMatcher(golexer.NewSignMatcher(Token_Comma, ":"))

	l.AddMatcher(golexer.NewIdentifierMatcher(Token_Identifier))

	l.AddMatcher(golexer.NewUnknownMatcher(Token_Unknown))

	l.Start(value)

	return &structParser{
		golexer.NewParser(l, value),
	}

}

func parseStruct(fd *model.FieldDescriptor, value string, fileD *model.FileDescriptor, node *model.Node) bool {

	p := newStructParser(value)

	// 检查字段有没有重复
	sfList := newStructFieldList()

	result := p.Run(fd, func(key, value string) bool {

		bnField := fd.Complex.FieldByValueAndMeta(key)
		if bnField == nil {

			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_FieldNotFound), key)

			return false
		}

		if sfList.Exists(bnField) {
			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_DuplicateFieldInCell), key)
			return false
		}

		sfList.Add(bnField, value)

		return true
	})

	if !result {
		return false
	}

	// 结构体中未填的字段如果是Default, 也要输出
	for _, structField := range fd.Complex.Fields {

		if sfList.Exists(structField) {
			continue
		}

		if structField.Meta.GetString("Default") != "" {
			sfList.Add(structField, structField.Meta.GetString("Default"))
		}

	}

	// 结构体输出是map顺序, 必须按照定义时的order进行排序, 否则在二进制中顺序是错的
	sfList.Sort()

	for i := 0; i < sfList.Len(); i++ {

		v := sfList.Get(i)

		// 添加类型节点
		fieldNode := node.AddKey(v.key)

		// 在类型节点下添加值节点
		_, ok := ConvertValue(v.key, v.value, fileD, fieldNode)

		if !ok {
			return false
		}

	}

	return true

}
