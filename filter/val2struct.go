package filter

import (
	"github.com/davyxu/tabtoy/proto/tool"
)

const (
	stateKey   = 1
	stateComma = 2
	stateValue = 3
)

func Value2Struct(meta *tool.FieldMeta, value string, callback func(string, string) bool) (isValue2Struct bool, hasError bool) {

	if meta == nil {
		return
	}

	if meta.String2Struct == false {
		return
	}

	lex := newLineLexer(value)

	parserState := stateKey
	var key string

	for {

		token, state := lex.Next()

		switch state {
		case lexerEOF:
			isValue2Struct = true
			return
		case lexerErr:
			hasError = true
			return
		case lexerToken:

			switch parserState {
			case stateKey:
				key = token
				parserState = stateComma
			case stateComma:
				if token != ":" {
					hasError = true
					log.Errorf("Unexpect symbol '%v' expect ':'", token)
					return
				}

				parserState = stateValue

			case stateValue:

				if !callback(key, token) {
					hasError = true
					return
				}

				parserState = stateKey
			}

		}

	}

	isValue2Struct = true
	return

}
