package model

import (
	"sort"

	"github.com/davyxu/golexer"
)

var builtinTag = map[string]bool{
	"MakeIndex":   true,
	"Alias":       true,
	"Default":     true,
	"ListSpliter": true,
	"RepeatCheck": true,
	"TableName":   true,
	"Package":     true,
	"OutputTag":   true,
}

func IsSystemTag(tag string) bool {
	_, ok := builtinTag[tag]
	return ok
}

type MetaInfo struct {
	*golexer.KVPair
}

func (self *MetaInfo) VisitUserMeta(callback func(string, interface{}) bool) {

	sortedKeys := make([]string, 0)

	for k, _ := range self.Raw() {

		if IsSystemTag(k) {
			continue
		}

		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		v, _ := self.Raw()[k]

		if !callback(k, v) {
			return
		}
	}

}

func NewMetaInfo() *MetaInfo {
	return &MetaInfo{
		KVPair: golexer.NewKVPair(),
	}
}
