package compiler

import (
	"github.com/davyxu/tabtoy/v4/model"
	"testing"
)

func TestMeta(t *testing.T) {

	header := model.NewHeaderField(0)
	parseMeta(header.TypeInfo, "")

	parseMeta(header.TypeInfo, " ")

	parseMeta(header.TypeInfo, "MakeIndex MakeIndex")

	if !header.TypeInfo.MakeIndex {
		t.FailNow()
	}

	parseMeta(header.TypeInfo, "ArraySpliter=|")
	if header.TypeInfo.ArraySplitter != "|" {
		t.FailNow()
	}
}
