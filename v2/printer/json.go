package printer

import (
	"encoding/json"

	"github.com/davyxu/tabtoy/v2/i18n"
	"github.com/davyxu/tabtoy/v2/model"
)

func init() { RegisterPrinter("json", &jsonPrinter{}) }

type jsonPrinter struct{}

func (self *jsonPrinter) Run(g *Globals) *Stream {
	v := make(map[string]interface{})
	for _, tab := range g.Tables {
		if !tab.LocalFD.MatchTag(".json") {
			log.Warnf("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		if !printTableJson(v, tab) {
			return nil
		}
	}

	buffer := NewStream()
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	buffer.WriteBytes(b)
	return buffer
}

func printTableJson(v map[string]interface{}, tab *model.Table) bool {
	fields := make([]interface{}, 0)
	for _, r := range tab.Recs {
		field := make(map[string]interface{})
		for _, node := range r.Nodes {
			if node.SugguestIgnore {
				continue
			}

			if node.Type != model.FieldType_Struct {
				if node.IsRepeated {
					subField := make([]interface{}, 0)
					for _, cnode := range node.Child {
						subField = append(subField, cnode.IValue)
					}
					field[node.Name] = subField
				} else {
					field[node.Name] = node.Child[0].IValue
				}
			} else {
				if node.IsRepeated {
					arrField := make([]interface{}, 0)
					for _, snode := range node.Child {
						if snode.SugguestIgnore {
							continue
						}
						vv := make(map[string]interface{})
						for _, cnode := range snode.Child {
							if cnode.SugguestIgnore {
								continue
							}
							vv[cnode.Name] = cnode.Child[0].IValue
						}
						arrField = append(arrField, vv)
					}
					field[node.Name] = arrField
				} else {
					subField := make(map[string]interface{})
					for _, snode := range node.Child {
						if snode.SugguestIgnore {
							continue
						}
						for _, cnode := range snode.Child {
							if cnode.SugguestIgnore {
								continue
							}
							subField[cnode.Name] = cnode.Child[0].IValue
						}
					}
					field[node.Name] = subField
				}
			}
		}
		fields = append(fields, field)
	}
	v[tab.LocalFD.Name] = fields
	return true
}
