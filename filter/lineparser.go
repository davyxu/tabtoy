package filter

type lineLexer struct {
	original string
	s        string
}

const (
	lexerEOF   = 0
	lexerToken = 1
	lexerErr   = 2
)

// Numbers and identifiers are matched by [-+._A-Za-z0-9]
func isIdentOrNumberChar(c byte) bool {
	switch {
	case 'A' <= c && c <= 'Z', 'a' <= c && c <= 'z':
		return true
	case '0' <= c && c <= '9':
		return true
	}
	switch c {
	case '-', '+', '.', '_':
		return true
	}
	return false
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	}
	return false
}

func (self *lineLexer) skipWhiteSpace() {

	i := 0
	for ; i < len(self.s); i++ {

		c := self.s[i]
		if isWhitespace(c) {
			continue
		}

		break
	}

	self.offset(i)
}

func (self *lineLexer) offset(p int) {
	self.s = self.s[p:len(self.s)]
}

func (self *lineLexer) Next() (string, int) {
	self.skipWhiteSpace()

	if len(self.s) == 0 {
		return "", lexerEOF
	}

	switch self.s[0] {
	case '"':

		var quotedString string
		var done bool

		i := 1
		for ; i < len(self.s); i++ {
			if self.s[i] == '"' {
				quotedString = self.s[1:i]
				done = true
				break
			}
		}

		if !done {
			log.Errorf("unfinished quoted string: %s", self.original)
			return "", lexerErr
		}
		self.offset(i + 1)

		return quotedString, lexerToken

	case ':':
		self.offset(1)

		return ":", lexerToken
	case '#':
		return "", lexerEOF

	}

	i := 0
	for i < len(self.s) && isIdentOrNumberChar(self.s[i]) {
		i++
	}

	if i == 0 {
		log.Errorf("unexpect char %c", self.s[0])
		return "", lexerErr
	}

	ret := self.s[0:i]

	self.offset(i)

	return ret, lexerToken
}

func newLineLexer(src string) *lineLexer {
	return &lineLexer{original: src,
		s: src,
	}
}
