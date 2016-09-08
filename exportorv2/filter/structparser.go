package filter

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/tabtoy/exportorv2/model"
)

type structParser struct {
	lexer *golexer.Lexer

	curr *golexer.Token
}

func (self *structParser) NextToken() {

	token, err := self.lexer.Read()

	if err != nil {
		panic(err)
	}

	if token == nil {
		panic("EOF")
		return
	}

	self.curr = token
}

func (self *structParser) TokenID() int {
	return self.curr.MatcherID()
}

func (self *structParser) TokenValue() string {
	return self.curr.Value()
}

func (self *structParser) Run(fd *model.FieldDefine, callback func(string, string) bool) (ok bool) {

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
			log.Errorf("field: '%s' parse error, %v", fd.Name, err)
		case string:
			if err.(string) == "EOF" {
				ok = true
			}

		}

	}()

	self.NextToken()

	for {

		if self.TokenID() != Token_Identifier {
			log.Errorf("expect field: '%s'", fd.Name)
			return false
		}

		key := self.TokenValue()

		self.NextToken()

		if self.TokenID() != Token_Comma {
			log.Errorf("%s need ':' split value", key)
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

// 自定义的token id
const (
	Token_Unknown = iota
	Token_WhiteSpace
	Token_LineEnd
	Token_UnixStyleComment
	Token_Identifier
	Token_Numeral
	Token_String
	Token_Comma
)

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
		lexer: l,
	}

}

func parseStruct(fd *model.FieldDefine, value string, typeset *model.BuildInTypeSet, node *model.Node) bool {

	p := newStructParser(value)

	return p.Run(fd, func(key, value string) bool {

		bnField := fd.BuildInType.FieldByValueAndMeta(key)
		if bnField == nil {
			log.Errorf("struct field not found: '%s'", key)
			return false
		}

		// 添加类型节点
		fieldNode := node.AddKey(bnField)

		// 在类型节点下添加值节点
		_, ok := ConvertValue(bnField, value, typeset, fieldNode)

		return ok
	})

}
