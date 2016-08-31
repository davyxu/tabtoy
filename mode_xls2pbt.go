package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"path/filepath"

	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/data"
	"github.com/davyxu/tabtoy/filter"
	"github.com/davyxu/tabtoy/printer"
	"github.com/davyxu/tabtoy/scanner"
)

///////////////////////////////////////////////
// mode: xls2pbt参数
///////////////////////////////////////////////

// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
var paramPbFile = flag.String("pb", "PB", "input protobuf binary descript file, export by protoc-gen-meta plugins")

// 输入电子表格文件
var paramXlsFile = flag.String("xls", "XLS", "input excel file, use ',' splited file list by multipy files")

// 输出文件夹
var paramOutDir = flag.String("outdir", "OUT_DIR", "output directory")

// 补丁文件
var paramPatch = flag.String("patch", "", "patch input files then output")

// 输出文件格式
var paramFormat = flag.String("fmt", "pbt", "output file format, support 'pbt', 'lua' ")

func runXls2PbtMode() bool {
	// 输入协议二进制描述文件,通过protoc配合github.com/davyxu/pbmeta/protoc-gen-meta插件导出
	// 创建描述文件池
	pool, err := pbmeta.CreatePoolByFile(*paramPbFile)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var patchFile *PatchFile

	if *paramPatch != "" {

		patchFile = NewPatchFile(pool)

		if !patchFile.Load(*paramPatch) {
			return false
		}

		log.Infof("patch file loaded: %s", *paramPatch)
	}

	return parallelWorker(flag.Args(), *paramPara, *paramOutDir, func(input, output string) bool {

		// 显示电子表格到导出文件

		sheetDataArray := exportSheetMsg(pool, input)

		if sheetDataArray == nil {
			return false
		}

		if patchFile != nil {
			patchFile.Patch(sheetDataArray)
		}

		return printFile(sheetDataArray, output)

	})

}

func setFieldValue(ri *scanner.RecordInfo, fieldName, value string) bool {

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

func getOutputExt() string {
	switch *paramFormat {
	case "json":
		return ".json"
	case "pbt":
		return ".pbt"
	case "lua":
		return ".lua"
	}

	return ""
}

func newWriter(buf *bytes.Buffer) printer.IWriter {
	switch *paramFormat {
	case "json":
		return printer.NewJsonWriter(buf)
	case "pbt":
		return printer.NewPBTWriter(buf)
	case "lua":
		return printer.NewLuaWriter(buf)
	}

	return nil
}

func printFile(sheetData []*sheetData, outputFile string) bool {

	log.Infof("		%s\n", filepath.Base(outputFile))

	var outBuff bytes.Buffer

	writer := newWriter(&outBuff)

	if writer == nil {
		log.Errorf("unknown output writer: %s\n", *paramFormat)
		return false
	}

	for _, sd := range sheetData {

		if !writer.PrintMessage(sd.msg) {
			return false
		}
	}

	// 创建输出文件
	file, err := os.Create(outputFile)
	if err != nil {
		log.Errorln(err.Error())
		return false
	}

	// 写入文件头

	file.WriteString(outBuff.String())

	file.Close()

	return true
}

func exportSheetMsg(pool *pbmeta.DescriptorPool, inputXls string) []*sheetData {

	// 显示要输出的文件
	log.Infof("%s\n", filepath.Base(inputXls))

	// 打开电子表格
	xlsFile := scanner.NewFile(pool)

	if !xlsFile.Open(inputXls) {
		return nil
	}

	sheetDataArray := make([]*sheetData, 0)

	repChecker := filter.NewRepeatValueChecker()

	// 遍历所有表格sheet
	for _, sheet := range xlsFile.Sheets {

		// 遍历表格的所有行/列
		sheetMsg, ok := sheet.IterateData(func(ri *scanner.RecordInfo) bool {

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
