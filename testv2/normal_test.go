package test

import (
	"io/ioutil"
	"testing"

	"github.com/davyxu/tabtoy/testv2/test"
	"github.com/golang/protobuf/proto"
)

func Test(t *testing.T) {

	data, err := ioutil.ReadFile("Actor.pbt")

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

	if file.Actor[1].Name != "葫芦\n娃" {
		t.Fatal("fail 1", file.Actor[1].Name)
	}

	if file.Actor[2].Name != "舒\"克\"" {
		t.Fatal("fail 2", file.Actor[2].Name)
	}

}
