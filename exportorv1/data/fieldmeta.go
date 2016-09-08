package data

import (
	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/golang/protobuf/proto"
)

// 获取一个字段的扩展信息
func GetFieldMeta(field interface{}) (*tool.FieldMetaV1, bool) {

	var metaStr string

	cm := field.(interface {
		ParseTaggedComment() []*pbmeta.TaggedComment
	})

	tc := cm.ParseTaggedComment()

	for _, c := range tc {

		if c.Name == "@" || c.Name == "tabtoy" {
			metaStr = c.Data
			break
		}

	}

	var meta tool.FieldMetaV1

	if err := proto.UnmarshalText(metaStr, &meta); err != nil {
		log.Errorf("parse field meta failed, [%s] %s", metaStr, err.Error())
		return nil, false
	}

	return &meta, true
}
