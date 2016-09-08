package exportor

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/exportorv1/data"
	"github.com/davyxu/tabtoy/exportorv1/filter"
	"github.com/davyxu/tabtoy/exportorv1/printer"
	"github.com/davyxu/tabtoy/util"
)

type Parameter struct {
	InputFileList []string
	PBFile        string
	PatchFile     string
	ParaMode      bool
	OutDir        string
	Format        string
}

func Run(param Parameter) bool {
	// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
	// 创建描述文件池
	pool, err := pbmeta.CreatePoolByFile(param.PBFile)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var patchFile *PatchFile

	if param.PatchFile != "" {

		patchFile = NewPatchFile(pool)

		if !patchFile.Load(param.PatchFile) {
			return false
		}

		log.Infof("patch file loaded: %s", param.PatchFile)
	}

	return util.ParallelWorker(param.InputFileList, param.ParaMode, func(input string) bool {

		var ext string
		var writer printer.IPrinter

		switch param.Format {
		case "pbt":
			ext = ".pbt"
			writer = printer.NewPBTWriter()
		case "json":
			ext = ".json"
			writer = printer.NewJsonWriter()

		case "lua":
			ext = ".lua"
			writer = printer.NewLuaWriter()

		default:
			log.Errorf("unknown format '%s'", param.Format)
			return false
		}

		// 使用指定的导出文件夹,并更换电子表格输入文件的后缀名为pbt作为输出文件
		outputFile := path.Join(param.OutDir, util.ChangeExtension(input, ext))

		// 显示电子表格到导出文件

		sheetDataArray := exportSheetMsg(pool, input)

		if sheetDataArray == nil {
			return false
		}

		if patchFile != nil {
			patchFile.Patch(sheetDataArray)
		}

		return printFile(sheetDataArray, outputFile, writer)

	})

}

func setFieldValue(ri *RecordInfo, fieldName, value string) bool {

	// 转换电子表格的原始值到msg可接受的值
	if afterValue, ok := filter.ValueConvetor(ri.FieldDesc, value); ok {

		fd := filter.FieldByNameWithMeta(ri.FieldMsg.Desc, fieldName)

		if fd == nil {
			log.Errorf("field not exist, field: '%s' value: '%s'", fieldName, value)
			return false
		}

		if data.DebuggingLevel >= 2 {
			log.Debugf("	%s|'%s'='%s'", ri.FieldMsg.Desc.Name(), fd.Name(), afterValue)
		}

		var ret bool

		// 多值
		if fd.IsRepeated() {
			ret = ri.FieldMsg.AddRepeatedValue(fd, afterValue)

		} else {

			// 单值
			ret = ri.FieldMsg.SetValue(fd, afterValue)
		}

		if !ret {
			log.Errorln("set value failed ", fd.Name(), afterValue)
		}

		return ret

	}

	log.Errorf("value convert error: '%s'='%s' field type: %s repteated: %v", fieldName, value, ri.FieldDesc.Type(), ri.FieldDesc.IsRepeated())
	return false

}

type sheetData struct {
	name string
	msg  *data.DynamicMessage // 对应XXFile
}

func printFile(sheetData []*sheetData, outputFile string, printer printer.IPrinter) bool {

	log.Infof("		%s\n", filepath.Base(outputFile))

	for _, sd := range sheetData {

		if !printer.PrintMessage(sd.msg) {
			return false
		}
	}

	return printer.WriteToFile(outputFile)

}

func exportSheetMsg(pool *pbmeta.DescriptorPool, inputXls string) []*sheetData {

	// 显示要输出的文件
	log.Infof("%s\n", filepath.Base(inputXls))

	// 打开电子表格
	xlsFile := NewFile(pool)

	if !xlsFile.Open(inputXls) {
		return nil
	}

	sheetDataArray := make([]*sheetData, 0)

	repChecker := filter.NewRepeatValueChecker()

	// 遍历所有表格sheet
	for _, sheet := range xlsFile.Sheets {

		// 遍历表格的所有行/列
		sheetMsg, ok := sheet.IterateData(func(ri *RecordInfo) bool {

			// 重复值检查
			if !repChecker.Check(ri.FieldMeta, ri.FieldDesc, ri.Value()) {
				return false
			}

			// 字符串转结构体
			v2sAffected, v2sHasErr := filter.Value2Struct(ri.FieldMeta, ri.Value(), ri.FieldDesc, func(key, value string) bool {

				return setFieldValue(ri, key, value)
			})

			if v2sHasErr {
				return false
			}

			if v2sAffected {
				return true
			}

			var setFieldValueHasErr bool

			// 分隔符切分值
			if filter.Value2List(ri.FieldMeta, ri.Value(), func(value string) {
				if !setFieldValue(ri, ri.FieldDesc.Name(), value) {
					setFieldValueHasErr = true
				}
			}) {
				return !setFieldValueHasErr
			}

			if setFieldValueHasErr {
				return false
			}

			return setFieldValue(ri, ri.FieldDesc.Name(), ri.Value())

		})

		if !ok {
			return nil
		}

		if sheetMsg != nil {
			// 显示导出Sheet时的名称
			log.Infof("	%s", sheet.Name)

			sheetDataArray = append(sheetDataArray, &sheetData{
				name: sheet.Name,
				msg:  sheetMsg,
			})
		}

	}

	return sheetDataArray
}
