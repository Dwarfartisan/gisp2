package gisp

import (
	"testing"
	"fmt"
	"reflect"
)

type lambda struct {
	index int
}

func (l lambda) Task(env Env, args ...interface{}) (Lisp, error){
	return TaskBox{func(env Env) (interface{}, error) {
		return l.index, nil
	}}, nil
}

func newLambda(idx int) lambda {
	return lambda{idx};
}

func TestLambda(t *testing.T){
	g := NewGisp(map[string]Toolbox{
		"axioms": Axiom,
		"props":  Propositions,
	})

	for idx :=0; idx< 10; idx++ {
		g.Defun(fmt.Sprintf("lambda%v", idx), newLambda(idx))
	}
	for idx :=0; idx< 10; idx++ {
		fname := fmt.Sprintf("lambda%v", idx);
		f, ok := g.Lookup(fname);
		if !ok {
			t.Errorf("%s not found", fname);
		}
		result, err := g.Eval(List{f});
		if err != nil {
			t.Errorf("%s call failed: %v", fname, err);
		}
		if !reflect.DeepEqual(result, idx){
			t.Errorf("expect %d but %v", idx, result);
		}
	}

}
