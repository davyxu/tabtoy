package exportorv2

import (
	"testing"
)

func TestParseFieldHeader(t *testing.T) {

	for {
		tp, r := ParseFieldTypeString("repeated float")
		if r != true || tp != FieldType_Float {
			t.Failed()

		}
		break
	}

	for {
		tp, r := ParseFieldTypeString("uint64")
		if r != false || tp != FieldType_UInt64 {
			t.Failed()
		}

		break
	}

}
