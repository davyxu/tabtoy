package filter

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
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
		golexer.NewParser(l),
	}

}

func parseStruct(fd *model.FieldDescriptor, value string, fileD *model.FileDescriptor, node *model.Node) bool {

	p := newStructParser(value)

	// 检查字段有没有重复
	fieldByFD := make(map[*model.FieldDescriptor]bool)

	return p.Run(fd, func(key, value string) bool {

		bnField := fd.Complex.FieldByValueAndMeta(key)
		if bnField == nil {

			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_FieldNotFound), key)

			return false
		}

		if _, ok := fieldByFD[bnField]; ok {
			log.Errorf("%s, '%s'", i18n.String(i18n.StructParser_DuplicateFieldInCell), key)
			return false
		}

		fieldByFD[bnField] = true

		// 添加类型节点
		fieldNode := node.AddKey(bnField)

		// 在类型节点下添加值节点
		_, ok := ConvertValue(bnField, value, fileD, fieldNode)

		return ok
	})

}
