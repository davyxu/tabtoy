package jsondata

// 报错行号+3
const templateText = `{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "{{.Version}}",	{{range $di, $tab := .Datas.AllTables}}
	"{{$tab.HeaderType}}":[ {{range $unusedrow,$row := $tab.DataIndexs}}
		{ {{range $col, $header := $tab.Headers}}"{{$header.TypeInfo.FieldName}}": {{WrapTabValue $ $tab $row $col}}{{GenJsonTailComma $col $tab.Headers}} {{end}}}{{GenJsonTailComma $row $tab.Rows}}{{end}} 
	]{{GenJsonTailComma $di $.Datas.AllTables}}{{end}}
}`
