package data

import (
	"strings"

	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/golang/protobuf/proto"
)

type IComment interface {
	TrailingComment() string

	LeadingComment() string
}

const headMarker = "@"

// 获取一个字段的扩展信息
func GetFieldMeta(cm IComment) *tool.FieldMeta {

	var metaStr string

	if strings.Index(cm.TrailingComment(), headMarker) == 0 {
		metaStr = strings.TrimSpace(cm.TrailingComment()[1:])
	} else if strings.Index(cm.LeadingComment(), headMarker) == 0 {
		metaStr = strings.TrimSpace(cm.LeadingComment()[1:])
	} else {
		return nil
	}

	var meta tool.FieldMeta

	if err := proto.UnmarshalText(metaStr, &meta); err != nil {
		log.Errorf("parse field meta failed, [%s] %s", metaStr, err.Error())
		return nil
	}

	return &meta
}
