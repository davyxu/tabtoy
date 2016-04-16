package printer

import (
	"testing"
)

func TestStrWrapper(t *testing.T) {
	t.Log(strEscape("a\"b\\nx"))
}
