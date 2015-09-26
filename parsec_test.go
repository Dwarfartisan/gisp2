package gisp

import (
	//"fmt"

	"reflect"
	"testing"

	p "github.com/Dwarfartisan/goparsec2"
)

func TestParsecBasic(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "p": Parsec})

	digit := p.Many1(p.Digit).Bind(p.ReturnString)
	data := "344932454094325"
	state := NewStringState(data)
	pcre, err := digit(state)
	if err != nil {
		t.Fatalf("expect \"%v\" pass test many1 digit but error:%v", data, err)
	}

	src := "(let ((st (p.state \"" + data + `")))
	(var data ((p.many1 p.digit) st))
	(p.s2str data))
	`
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("expect \"%v\" pass gisp many1 digit but error:%v", src, err)
	}
	t.Logf("from gisp: %v", gre)
	t.Logf("from parsec: %v", pcre)
	if !reflect.DeepEqual(pcre, gre) {
		t.Fatalf("expect got \"%v\" from gisp equal \"%v\" from parsex", gre, pcre)
	}
}

func TestParsecRune(t *testing.T) {
	g := NewGispWith(
		map[string]Toolbox{
			"axiom": Axiom, "props": Propositions, "time": Time},
		map[string]Toolbox{"time": Time, "p": Parsec})
	//data := "Here is a Rune : 'a' and a is't a rune. It is a word in sentence."
	data := "'a' and a is't a rune. It is a word in sentence."
	state := NewStringState(data)
	pre, err := p.Between(p.Chr('\''), p.Chr('\''), p.AsRune)(state)
	if err != nil {
		t.Fatalf("expect found rune expr from \"%v\" but error:%v", data, err)
	}
	src := `
	(let ((st (p.state "` + data + `")))
		((p.between (p.rune '\'') (p.rune '\'') p.one) st))
	`
	gre, err := g.Parse(src)
	if err != nil {
		t.Fatalf("expect \"%v\" pass gisp '<rune>' but error:%v", src, err)
	}
	t.Logf("from gisp: %v", gre)
	t.Logf("from parsec: %v", pre)
	if !reflect.DeepEqual(pre, gre) {
		t.Fatalf("expect got \"%v\" from gisp equal \"%v\" from parsec", gre, pre)
	}
}
