package gisp

import (
	"testing"

	p "github.com/Dwarfartisan/goparsec2"
)

func TestIntParser0(t *testing.T) {
	data := "12"
	st := p.BasicStateFromText(data)
	o, err := IntParser(&st)
	if err != nil {
		t.Fatalf("expect a Int but error %v", err)
	}
	if i, ok := o.(Int); ok {
		if i != Int(12) {
			t.Fatalf("expect a Int 12 but %v", i)
		}
	} else {
		t.Fatalf("expect Int but %v", o)
	}
}

func TestIntParser1(t *testing.T) {
	data := "i234"
	st := p.BasicStateFromText(data)
	o, err := IntParser(&st)
	if err == nil {
		t.Fatalf("expect a Int parse error but got %v", o)
	}
}

func TestIntParser2(t *testing.T) {
	data := ".234"
	st := p.BasicStateFromText(data)
	o, err := IntParser(&st)
	if err == nil {
		t.Fatalf("expect a Float parse error but got %v", o)
	}
}

func TestIntParser3(t *testing.T) {
	data := "3.14"
	st := p.BasicStateFromText(data)
	o, err := p.M(IntParser).Then(p.EOF)(&st)
	if err == nil {
		t.Fatalf("expect a Float parse error but got %v", o)
	}
}
