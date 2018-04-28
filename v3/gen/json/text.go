package json

// 报错行号+3
const templateText = `{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "{{.Version}}",	{{range $di, $tab := .Datas}}
	"{{$tab.Name}}":[ {{range $row,$rowData := $tab.Rows}}
		{ {{range $col, $headType := $tab.HeaderFields}}"{{$headType.FieldName}}": {{WrapTabValue $ $tab $row $col}}{{GenJsonTailComma $col $tab.HeaderFields}} {{end}}}{{GenJsonTailComma $row $tab.Rows}}{{end}} 
	]{{GenJsonTailComma $di $.Datas}}{{end}}
}`
