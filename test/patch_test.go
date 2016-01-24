package test

import (
	"github.com/davyxu/tabtoy/test/test"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"os"
	"testing"
)

func TestPatch(t *testing.T) {

	f, err := os.Open("Actor.pbt")
	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	data, err := ioutil.ReadAll(f)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	var file test.ActorFile

	err = proto.UnmarshalText(string(data), &file)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	if file.Actor[0].Name != "A" {
		t.Fatalf("fail A")
	}

}
