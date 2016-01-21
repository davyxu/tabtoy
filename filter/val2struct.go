package filter

import (
	"github.com/davyxu/tabtoy/proto/tool"
)

const (
	stateKey   = 1
	stateComma = 2
	stateValue = 3
)

func Value2Struct(meta *tool.FieldMeta, value string, callback func(string, string)) bool {

	if meta == nil {
		return false
	}

	if meta.String2Struct == false {
		return false
	}

	lex := newLineLexer(value)

	parserState := stateKey
	var key string

	for {

		token, state := lex.Next()

		if state == lexerErr ||
			state == lexerEOF {
			break
		}

		if state == lexerToken {

			switch parserState {
			case stateKey:
				key = token
				parserState = stateComma
			case stateComma:
				if token != ":" {
					log.Errorf("Unexpect symbol '%v' expect ':'", token)
					return true
				}

				parserState = stateValue

			case stateValue:

				callback(key, token)

				parserState = stateKey
			}

		}

	}

	return true

}
