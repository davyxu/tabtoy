package exportor

import (
	"github.com/davyxu/pbmeta"
	pbprotos "github.com/davyxu/pbmeta/proto"
	"github.com/davyxu/tabtoy/data"
)

type PatchFile struct {
	pool *pbmeta.DescriptorPool

	sheetDataArray []*sheetData
}

func (self *PatchFile) Load(filename string) bool {

	self.sheetDataArray = exportSheetMsg(self.pool, filename)

	if self.sheetDataArray == nil {
		return false
	}

	return true
}

func (self *PatchFile) Patch(input []*sheetData) bool {

	if self.sheetDataArray == nil {
		return false
	}

	if len(self.sheetDataArray) != len(input) {
		log.Errorf("patch failed, sheet count do not match")
		return false
	}

	for index, sheetData := range self.sheetDataArray {

		inputData := input[index]
		if inputData.name != sheetData.name {
			log.Errorf("patch failed, sheet name not match")
			return false
		}

		if !patchMsg(inputData.msg, sheetData.msg, 0, inputData.name) {
			return false
		}
	}

	return true
}

func patchMsg(input, patch *data.DynamicMessage, indent int, sheetName string) bool {

	return patch.IterateFieldDesc(func(fd *pbmeta.FieldDescriptor) bool {

		if fd.Type() == pbprotos.FieldDescriptorProto_TYPE_MESSAGE {

			if fd.IsRepeated() {

				// 顶级的消息肯定是repeated的, 因此直接进入
				if indent == 0 {

					nextInputMsgArray := input.GetRepeatedMessage(fd)
					nextPatchMsgArray := patch.GetRepeatedMessage(fd)

					if len(nextInputMsgArray) != len(nextPatchMsgArray) {
						log.Errorf("patch rows diff from input msg")
						return false
					}

					// 相同的行数才可以进去继续迭代
					for index, nextInputMsg := range nextInputMsgArray {

						patchMsg(nextInputMsg, nextPatchMsgArray[index], indent+1, sheetName)
					}

				} else {

					// 直接替换值数组
					pathMsgArray := patch.GetRepeatedMessage(fd)

					input.ClearFieldValue(fd)
					for _, msg := range pathMsgArray {
						input.AddRepeatedMessage(fd, msg)
					}

				}

			} else {

				nextInputMsg := input.GetMessage(fd)
				nextPatchMsg := patch.GetMessage(fd)

				patchMsg(nextInputMsg, nextPatchMsg, indent+1, sheetName)
			}

		} else {

			if fd.IsRepeated() {

				log.Infof("patch repeated field: '%s' = '%s'", fd.Name(), patch.GetRepeatedValue(fd))

				// 直接用patch的覆盖多值
				input.ClearFieldValue(fd)
				for _, v := range patch.GetRepeatedValue(fd) {

					input.AddRepeatedValue(fd, v)
				}

			} else {

				v, _ := patch.GetValue(fd)

				log.Infof("patch field: '%s' = '%s'", fd.Name(), v)

				// 单值设置
				input.SetValue(fd, v)

			}

		}

		return true
	})

	return true
}

func NewPatchFile(pool *pbmeta.DescriptorPool) *PatchFile {

	return &PatchFile{
		pool: pool,
	}

}
