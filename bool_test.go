package gisp

import (
	"testing"

	p "github.com/Dwarfartisan/goparsec2"
)

func TestBoolParse0(t *testing.T) {
	data := "true"
	st := p.BasicStateFromText(data)
	o, err := BoolParser(&st)
	if err != nil {
		t.Fatalf("expect bool but error %v", err)
	}
	if b, ok := o.(Bool); ok {
		if !b {
			t.Fatalf("expect bool true but %v", b)
		}
	} else {
		t.Fatalf("excpet bool but %v", o)
	}
}

func TestBoolParse1(t *testing.T) {
	data := "false"
	st := p.BasicStateFromText(data)
	o, err := BoolParser(&st)
	if err != nil {
		t.Fatalf("expect bool but error %v", err)
	}
	if b, ok := o.(Bool); ok {
		if b {
			t.Fatalf("expect bool false but %v", b)
		}
	} else {
		t.Fatalf("excpet bool but %v", o)
	}
}
