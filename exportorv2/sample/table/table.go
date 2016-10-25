package table

// table的索引入口函数

var (
	indexEntryByName = make(map[string]func(interface{}))
)

func RegisterIndexEntry(name string, callback func(interface{})) {

	if _, ok := indexEntryByName[name]; ok {
		panic("duplicate table index entry")
	}

	indexEntryByName[name] = callback
}

func MakeIndex(content interface{}) {

	for _, v := range indexEntryByName {
		v(content)
	}

}
