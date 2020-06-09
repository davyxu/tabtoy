package jsondata

// 报错行号+3
const templateText = `{
	"@Tool": "github.com/davyxu/tabtoy",
	"@Version": "{{.Version}}",	{{range $di, $tab := .Datas.AllTables}}
	"{{$tab.HeaderType}}":[ {{range $unusedrow,$row := $tab.DataRowIndex}}{{$headers := $.Types.AllFieldByName $tab.OriginalHeaderType}}
		{ {{range $col, $header := $headers}}"{{$header.FieldName}}": {{WrapTabValue $ $tab $headers $row $col}}{{GenJsonTailComma $col $headers}} {{end}}}{{GenJsonTailComma $row $tab.Rows}}{{end}} 
	]{{GenJsonTailComma $di $.Datas.AllTables}}{{end}}
}`
