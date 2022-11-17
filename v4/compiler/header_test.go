package compiler

import (
	"github.com/davyxu/tabtoy/v4/model"
	"testing"
)

func TestMeta(t *testing.T) {

	header := &model.HeaderField{
		Cell:     &model.Cell{},
		TypeInfo: &model.DataField{},
	}
	parseMeta(header, ";")

	parseMeta(header, ";;")

	parseMeta(header, "MakeIndex;;MakeIndex")

	if !header.TypeInfo.MakeIndex {
		t.FailNow()
	}

	parseMeta(header, ";Spliter=|;")
	if header.TypeInfo.ArraySplitter != "|" {
		t.FailNow()
	}
}
