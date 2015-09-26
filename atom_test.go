package gisp

import (
	"reflect"
	"testing"

	p "github.com/Dwarfartisan/goparsec2"
)

func TestAtomParse0(t *testing.T) {
	data := "x"
	state := p.BasicStateFromText(data)
	a, err := AtomParser(&state)
	if err == nil {
		test := Atom{"x", Type{ANY, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("expect Atom{\"x\", ANY} but %v", a)
		}
	} else {
		t.Fatalf("expect Atom{\"x\", ANY} but %v", err)
	}
}

func TestAtomParse1(t *testing.T) {
	data := "x::atom"
	state := p.BasicStateFromText(data)
	a, err := AtomParser(&state)
	if err == nil {
		test := Atom{"x", Type{ATOM, false}}
		d := a.(Atom)
		if !reflect.DeepEqual(test, d) {
			t.Fatalf("expect Atom{\"x\", ATOM} but {Name:%v, Type:%v}", d.Name, d.Type)
		}
	} else {
		t.Fatalf("expect Atom{\"x\", ATOM} but error %v", err)
	}
}

func TestAtomParse2(t *testing.T) {
	data := "x::any"
	state := p.BasicStateFromText(data)
	a, err := AtomParser(&state)
	if err == nil {
		test := Atom{"x", Type{ANY, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("expect Atom{\"x\", ANY} but %v", a)
		}
	} else {
		t.Fatalf("expect Atom{\"x\", ANY} but %v", err)
	}
}

func TestAtomParse3(t *testing.T) {
	data := "x::int"
	state := p.BasicStateFromText(data)
	a, err := AtomParser(&state)
	if err == nil {
		test := Atom{"x", Type{INT, false}}
		if !reflect.DeepEqual(test, a) {
			t.Fatalf("expect Atom{\"x\", INT} but %v", a)
		}
	} else {
		t.Fatalf("expect Atom{\"x\", INT} but %v", err)
	}
}
