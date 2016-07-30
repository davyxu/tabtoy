package test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/davyxu/tabtoy/test/test"
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

	outdata, err := json.MarshalIndent(&file, "", " ")
	ioutil.WriteFile("Actor_std.json", outdata, 777)

	testJson(t, &file)

}

func testJson(t *testing.T, rightJson *test.ActorFile) {

	data, err := ioutil.ReadFile("Actor.json")

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	var file test.ActorFile

	err = json.Unmarshal(data, &file)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	testdata, err := proto.Marshal(&file)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	rightdata, err := proto.Marshal(rightJson)

	if err != nil {
		t.Fatalf(err.Error())
		return
	}

	if len(testdata) != len(rightdata) {
		t.Fatal("json not ok")
	}

	for i, d := range testdata {
		if rightdata[i] != d {
			t.Fatal("json not ok")
		}
	}

}
