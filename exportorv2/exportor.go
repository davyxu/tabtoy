package exportorv2

import (
	"bytes"
	"path"
	"path/filepath"

	"github.com/davyxu/tabtoy/util"
)

type Parameter struct {
	InputFileList []string
	ParaMode      bool
	OutDir        string
	Format        string
}

func printIndent(indent int) string {

	var buf bytes.Buffer
	for i := 0; i < indent*10; i++ {
		buf.WriteString(" ")
	}

	return buf.String()
}

func Run(param Parameter) bool {
	return util.ParallelWorker(param.InputFileList, param.ParaMode, param.OutDir, func(input string) bool {

		//	 显示电子表格到导出文件

		file := NewFile()

		tab := file.Export(input)
		if tab == nil {
			return false
		}

		if !tab.Print(file.Name) {
			return false
		}

		var ext string

		switch param.Format {
		case "pbt":
			ext = ".pbt"
		case "json":
			ext = ".json"

		case "lua":
			ext = ".lua"

		default:
			log.Errorf("unknown format '%s'", param.Format)
			return false
		}

		// 使用指定的导出文件夹,并更换电子表格输入文件的后缀名为pbt作为输出文件
		outputFile := path.Join(param.OutDir, util.ChangeExtension(input, ext))

		log.Infof("%s%s\n", printIndent(2), filepath.Base(outputFile))
		return tab.WriteToFile(outputFile)

	})

}
