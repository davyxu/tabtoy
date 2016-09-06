package util

import (
	"testing"
)

func TestA1toR1C1(t *testing.T) {

	if index2Alphabet(1) != "A" {
		t.Fatal("A")
	}

	if index2Alphabet(25) != "Y" {
		t.Fatal("Y")
	}

	if index2Alphabet(26) != "Z" {
		t.Fatal("Z")
	}

	if index2Alphabet(135) != "EE" {
		t.Fatal("EE")
	}

	if index2Alphabet(675) != "YY" {
		t.Fatal("YY")
	}

	if index2Alphabet(676) != "YZ" {
		t.Fatal("YZ")
	}

	if index2Alphabet(677) != "ZA" {
		t.Fatal("ZA")
	}

	if index2Alphabet(702) != "ZZ" {
		t.Fatal("ZZ")
	}

	if index2Alphabet(704) != "AAB" {
		t.Fatal("AAB")
	}

}
