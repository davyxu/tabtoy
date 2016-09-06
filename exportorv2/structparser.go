package exportorv2

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/pbmeta"
)

type structParser struct {
	lexer *golexer.Lexer
	fd    *pbmeta.FieldDescriptor

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

// 自定义的token id
const (
	Token_Unknown = iota
	Token_WhiteSpace
	Token_LineEnd
	Token_UnixStyleComment
	Token_Identifier
	Token_Comma
)

func newLineParser(fd *pbmeta.FieldDescriptor, value string) *structParser {
	l := golexer.NewLexer()

	//	l.AddMatcher(golexer.NewNumeralMatcher(Token_Numeral))
	//	l.AddMatcher(golexer.NewStringMatcher(Token_String))

	l.AddIgnoreMatcher(golexer.NewWhiteSpaceMatcher(Token_WhiteSpace))
	l.AddIgnoreMatcher(golexer.NewLineEndMatcher(Token_LineEnd))
	l.AddIgnoreMatcher(golexer.NewUnixStyleCommentMatcher(Token_UnixStyleComment))

	//	l.AddMatcher(golexer.NewSignMatcher(Token_True, "true"))
	//	l.AddMatcher(golexer.NewSignMatcher(Token_False, "false"))
	l.AddMatcher(golexer.NewSignMatcher(Token_Comma, ":"))

	l.AddMatcher(golexer.NewIdentifierMatcher(Token_Identifier))

	l.AddMatcher(golexer.NewUnknownMatcher(Token_Unknown))

	l.Start(value)

	return &structParser{
		lexer: l,
		fd:    fd,
	}
}

func parseStructString() {

}
