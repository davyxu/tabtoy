package model

import (
	"fmt"
	"strings"
)

const (
	ActionNoGenFieldJson     = "nogenfield_json"
	ActionNoGenFieldJsonDir  = "nogenfield_jsondir"
	ActionNoGenFieldBinary   = "nogenfield_binary"
	ActionNoGenFieldPbBinary = "nogenfield_pbbin"
	ActionNoGennFieldLua  	 = "nogenfield_lua"
	ActionNoGenTable         = "nogentab"
)

// 用tag选中目标, 做action
type TagAction struct {
	Verb string
	Tags []string
}

// action1:tag1+tag2|action2:tag1+tag3
func ParseTagAction(script string) (ret []TagAction, err error) {

	for _, as := range strings.Split(script, "|") {
		actionPairs := strings.Split(as, ":")

		var ta TagAction

		switch len(actionPairs) {
		case 2:
			ta.Verb = actionPairs[0]
			ta.Tags = strings.Split(actionPairs[1], "+")
			ret = append(ret, ta)
		default:
			err = fmt.Errorf("invalid action format")
			return
		}

	}

	return
}

func (self *Globals) CanDoAction(action string, obj interface{}) bool {

	for _, ta := range self.TagActions {
		if ta.Verb == action {
			for _, tag := range ta.Tags {
				switch v := obj.(type) {
				case *HeaderField:
					if v.TypeInfo.ContainTag(tag) {
						return true
					}
				case *TypeDefine:
					if v.ContainTag(tag) {
						return true
					}
				case *IndexDefine:
					if v.ContainTag(tag) {
						return true
					}
				}
			}
		}
	}

	return false
}
