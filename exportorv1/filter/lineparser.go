package filter

import (
	"github.com/davyxu/golexer"
	"github.com/davyxu/pbmeta"
)

type lineParser struct {
	lexer *golexer.Lexer
	fd    *pbmeta.FieldDescriptor

	curr *golexer.Token
}

func (self *lineParser) NextToken() {

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

func (self *lineParser) TokenID() int {
	return self.curr.MatcherID()
}

func (self *lineParser) TokenValue() string {
	return self.curr.Value()
}

func newLineParser(fd *pbmeta.FieldDescriptor, value string) *lineParser {
	l := golexer.NewLexer()

	l.AddMatcher(golexer.NewNumeralMatcher(Token_Numeral))
	l.AddMatcher(golexer.NewStringMatcher(Token_String))

	l.AddIgnoreMatcher(golexer.NewWhiteSpaceMatcher(Token_WhiteSpace))
	l.AddIgnoreMatcher(golexer.NewLineEndMatcher(Token_LineEnd))
	l.AddIgnoreMatcher(golexer.NewUnixStyleCommentMatcher(Token_UnixStyleComment))

	l.AddMatcher(golexer.NewKeywordMatcher(Token_True, "true"))
	l.AddMatcher(golexer.NewKeywordMatcher(Token_False, "false"))
	l.AddMatcher(golexer.NewSignMatcher(Token_Comma, ":"))

	l.AddMatcher(golexer.NewIdentifierMatcher(Token_Identifier))

	l.AddMatcher(golexer.NewUnknownMatcher(Token_Unknown))

	l.Start(value)

	return &lineParser{
		lexer: l,
		fd:    fd,
	}
}
