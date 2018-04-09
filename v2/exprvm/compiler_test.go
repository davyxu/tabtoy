package exprvm

import (
	"testing"
)

func TestCompiler(t *testing.T) {

	code := `-2+1`

	ck, err := Compile(code)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(ck)

	vm := NewMachine()
	vm.Run(ck)

	t.Log(vm.DataStack.String())
}
