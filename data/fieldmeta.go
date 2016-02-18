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

// 此样式无法支持多个tag, 即将被废除
const oldStyleMarker = "@"

const newStyleMarker = "[tabtoy]"

func findMetaString(marker string, cm IComment) string {
	if pos := strings.Index(cm.TrailingComment(), marker); pos == 0 {
		return strings.TrimSpace(cm.TrailingComment()[len(marker):])
	} else if strings.Index(cm.LeadingComment(), marker) == 0 {
		return strings.TrimSpace(cm.LeadingComment()[len(marker):])
	}

	return ""
}

// 获取一个字段的扩展信息
func GetFieldMeta(cm IComment) *tool.FieldMeta {

	var metaStr string

	// 优先找新样式, 再兼容老样式
	if metaStr = findMetaString(newStyleMarker, cm); metaStr == "" {
		if metaStr = findMetaString(oldStyleMarker, cm); metaStr == "" {
			return nil
		}
	}

	var meta tool.FieldMeta

	if err := proto.UnmarshalText(metaStr, &meta); err != nil {
		log.Errorf("parse field meta failed, [%s] %s", metaStr, err.Error())
		return nil
	}

	return &meta
}
