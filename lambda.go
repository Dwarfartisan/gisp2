package gisp

import (
	"fmt"

	p "github.com/Dwarfartisan/goparsec2"
)

// Lambda 实现基本的 Lambda 行为
type Lambda struct {
	Meta    map[string]interface{}
	Content List
}

// DeclareLambda 构造 Lambda 表达式 (lambda (args...) body)
func DeclareLambda(env Env, args List, lisps ...interface{}) (*Lambda, error) {
	ret := Lambda{map[string]interface{}{
		"category": "lambda",
		"local":    map[string]interface{}{},
	}, List{}}
	ret.prepareArgs(args)
	prepare := map[string]bool{}
	for _, lisp := range lisps {
		err := ret.prepare(env, prepare, lisp)
		if err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// LambdaExpr 生成一个封装后的 Lambda 表达式
func LambdaExpr(env Env, args ...interface{}) (Tasker, error) {
	st := p.NewBasicState(args)
	_, err := TypeAs(LIST)(&st)
	if err != nil {
		return nil, st.Trap("Lambda Args Error: expect args list but error: %v", err)
	}
	lptr, err := DeclareLambda(env, args[0].(List), args[1:]...)
	if err != nil {
		return nil, fmt.Errorf("Lambda Args Error: expect lambda tasker but error: %v", err)
	}
	return Q(lptr).Eval, nil
}

func (lambda *Lambda) prepareArgs(args List) {
	l := len(args)
	formals := make(List, len(args))
	if l == 0 {
		lambda.Meta["parameters parsex"] = []Var{}
		return
	}
	lidx := l - 1
	last := args[lidx].(Atom)
	// variadic function args formal as (last[::Type] ... )
	isVariadic := false
	if last.Name == "..." && len(args) > 1 {
		isVariadic = true
	}
	lambda.Meta["is variadic"] = isVariadic
	ps := make([]p.P, l+1)
	for idx, arg := range args[:lidx] {
		ps[idx] = argParser(arg.(Atom))
		formals[idx] = arg
	}
	if isVariadic {
		varArg := args[l-2].(Atom)
		ps[lidx] = p.Many(argParser(last))
		larg := Atom{varArg.Name, varArg.Type}
		formals[lidx] = larg
	} else {
		ps[lidx] = argParser(last)
		formals[lidx] = last
	}
	ps[l] = p.EOF
	lambda.Meta["formal parameters"] = formals
	lambda.Meta["parameter parsexs"] = ps
}

func (lambda *Lambda) prepare(env Env, prepare map[string]bool, content interface{}) error {
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}
	var err error
	switch lisp := content.(type) {
	case Atom:
		err = lambda.prepareAtom(env, next, lisp)
		return err
	case List:
		err = lambda.prepareList(env, next, lisp)
	}
	if err == nil {
		lambda.Content = append(lambda.Content, content)
	}
	return err
}

func (lambda Lambda) prepareAtom(env Env, prepare map[string]bool, one Atom) error {
	if _, ok := prepare[one.Name]; ok {
		return nil
	}
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}

	for _, arg := range lambda.Meta["formal parameters"].(List) {
		if (arg.(Atom)).Name == one.Name {
			return nil
		}
	}
	if _, ok := prepare[one.Name]; !ok {
		if v, ok := env.Lookup(one.Name); ok {
			local := (lambda.Meta["local"]).(map[string]interface{})
			local[one.Name] = v
		} else {
			return fmt.Errorf("%s not found", one.Name)
		}
	}
	return nil
}

func (lambda Lambda) prepareList(env Env, prepare map[string]bool, content List) error {
	next := map[string]bool{}
	for key := range prepare {
		next[key] = true
	}
	var err error
	fun := content[0].(Atom)
	switch fun.Name {
	case "var":
		name := content[1].(string)
		if err != nil {
			return err
		}
		next[name] = true
	case "lambda":
		args := content[1].(List)
		for _, a := range args {
			arg := a.(Atom)
			next[arg.Name] = true
		}
	case "let":
		for _, def := range content[1].(List) {
			define := def.(List)
			name := define[0].(string)
			next[name] = true
		}
	}
	for _, l := range content {
		switch lisp := l.(type) {
		case List:
			err = lambda.prepareList(env, next, lisp)
		case Atom:
			err = lambda.prepareAtom(env, next, lisp)
		}
	}
	return err
}

// TypeSign 生成反射类型签名
func (lambda Lambda) TypeSign() []Type {
	formals := lambda.Meta["formal parameters"].(List)
	types := make([]Type, len(formals))
	for idx, formal := range formals {
		types[idx] = formal.(Atom).Type
	}
	return types
}

// MatchArgsSign 校验参数是否匹配
func (lambda Lambda) MatchArgsSign(env Env, args ...interface{}) (interface{}, error) {
	params := make([]interface{}, len(args))
	for idx, arg := range args {
		param, err := Eval(env, arg)
		if err != nil {
			return nil, err
		}
		params[idx] = param
	}
	pxs := lambda.Meta["parameter parsexs"].([]p.P)
	st := p.NewBasicState(params)
	return p.UnionAll(pxs...)(&st)
}

// Task create a lambda s-Expr can be eval
func (lambda Lambda) Task(env Env, args ...interface{}) (Lisp, error) {
	meta := map[string]interface{}{}
	for k, v := range lambda.Meta {
		meta[k] = v
	}
	actuals, err := lambda.MatchArgsSign(env, args...)
	if err != nil {
		return Nil{}, err
	}
	meta["actual parameters"] = actuals
	meta["my"] = map[string]Var{}
	l := len(lambda.Content)
	content := make([]interface{}, l)
	for idx, data := range lambda.Content {
		content[idx] = data
	}
	return &Task{meta, content}, nil
}
