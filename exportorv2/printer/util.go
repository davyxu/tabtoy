package printer

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type FilePrinter struct {
	buf bytes.Buffer
}

func (self *FilePrinter) Data() []byte {
	return self.buf.Bytes()
}

func (self *FilePrinter) Printf(format string, args ...interface{}) {
	self.buf.WriteString(fmt.Sprintf(format, args...))
}

func (self *FilePrinter) Write(outfile string) bool {
	err := ioutil.WriteFile(outfile, self.buf.Bytes(), 0666)
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	return true
}
