package gisp

import (
	"fmt"
	//	p "github.com/Dwarfartisan/goparsec/parsex"
	"reflect"

	p "github.com/Dwarfartisan/goparsec2"
)

// Parsec 包为 gisp 解释器提供 parsec 解析工具
var Parsec = Toolkit{
	Meta: map[string]interface{}{
		"name":     "parsec",
		"category": "package",
	},
	Content: map[string]interface{}{
		"state": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Arg Error:expect args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch data := param.(type) {
			case string:
				return Q(NewStringState(data)), nil
			case List:
				return Q(p.NewBasicState(data)), nil
			default:
				return nil, fmt.Errorf("Parsex Error: expect create a state from a string or List but %v", data)
			}
		},
		"s2str": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Slice to string Arg Error:expect args has 1 arg.")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var (
				slice []interface{}
				ok    bool
			)
			if slice, ok = param.([]interface{}); !ok {
				return nil, fmt.Errorf("s2str Arg Error:expect 1 []interface{} arg")
			}
			return Q(p.ToString(slice)), nil
		},
		"one": ParsecBox(p.One),
		"eq": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Equal Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Eq(param)), nil
		},
		"str": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("One Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Str(param.(string))), nil
		},
		"rune": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Rune Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Chr(rune(param.(Rune)))), nil
		},
		"asint":   ParsecBox(p.AsInt),
		"asfloat": ParsecBox(p.AsFloat64),
		"asstr":   ParsecBox(p.AsString),
		"string": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("string Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var str string
			var ok bool
			if str, ok = param.(string); !ok {
				return nil, fmt.Errorf("string Arg Error:expect 1 string arg")
			}
			return ParsecBox(p.Str(str)), nil
		},
		"digit": ParsecBox(p.Digit),
		"int": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("int Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var i Int
			var ok bool
			if i, ok = param.(Int); !ok {
				return nil, fmt.Errorf("int Arg Error:expect 1 string arg")
			}
			return ParsecBox(func(st p.State) (interface{}, error) {
				data, err := p.Int(st)
				if err != nil {
					return nil, st.Trap("gisp parsex error:expect a int but error: %v", err)
				}
				if Int(data.(int)) != i {
					return nil, st.Trap("gisp parsex error:expect a Int but %v", data)
				}
				return data, nil
			}), nil
		},
		"float": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("float Arg Error:expect args has 1 arg")
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var f Float
			var ok bool
			if f, ok = param.(Float); !ok {
				return nil, fmt.Errorf("float Arg Error:expect 1 string arg")
			}
			return ParsecBox(func(st p.State) (interface{}, error) {
				data, err := p.AsFloat64(st)
				if err != nil {
					return nil, st.Trap("gisp parsex error:expect a float but error: %v", err)
				}
				if Float(data.(float64)) != f {
					return nil, st.Trap("gisp parsex error:expect a Float but %v", data)
				}
				return data, nil
			}), nil
		},
		"eof":    ParsecBox(p.EOF),
		"nil":    ParsecBox(p.Nil),
		"atimex": ParsecBox(TimeValue),
		"try": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Try Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Try(parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Try Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}

		},
		"either": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((p.P)(nil))
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Choice(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"choice": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((p.P)(nil))
			params, err := GetArgs(env, p.ManyTil(TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			parsers := make([]p.P, len(params))
			for idx, prs := range params {
				if parser, ok := prs.(Parsecer); ok {
					parsers[idx] = parser.Parser
				}
				return nil, fmt.Errorf("Choice Args Error:expect parsec parsers but %v is %v",
					prs, reflect.TypeOf(prs))
			}
			return ParsecBox(p.Choice(parsers...)), nil
		},
		"return": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Return Error: only accept one parsec parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Return(param)), nil
		},
		"option": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("Parsex Parser Option Error: only accept two parsex parser as arg but %v", args)
			}
			data, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			param, err := Eval(env, args[1])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Option(data, parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Many Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many1": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Many1 Erroparserr: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Many1(parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Many1 Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"many": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Many Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Many(parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Many Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"failed": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Failed Error: only accept one string as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var str string
			var ok bool
			if str, ok = param.(string); !ok {
				return nil, fmt.Errorf("Failed Arg Error:expect 1 string arg.")
			}
			return ParsecBox(p.Fail(str)), nil
		},
		"oneof": func(env Env, args ...interface{}) (Lisp, error) {
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.OneOf(params...)), nil
		},
		"noneof": func(env Env, args ...interface{}) (Lisp, error) {
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.NoneOf(params)), nil
		},
		"between": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.Between(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser, params[2].(Parsecer).Parser)), nil
		},
		"bind": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("Bind Args Error:expect 2 args.")
			}
			prs, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			var parser Parsecer
			var ok bool
			if parser, ok = prs.(Parsecer); !ok {
				return nil, fmt.Errorf("Bind Args Error:expect first arg is a Parsecer.")
			}
			f, err := Eval(env, args[1])
			if err != nil {
				return nil, err
			}
			switch fun := f.(type) {
			case func(interface{}) p.P:
				return ParsecBox(parser.Parser.Bind(fun)), nil
			case Functor:
				return ParsecBox(parser.Parser.Bind(func(x interface{}) p.P {
					tasker, err := fun.Task(env, x)
					if err != nil {
						return func(st p.State) (interface{}, error) {
							return nil, err
						}
					}
					pr, err := tasker.Eval(env)
					if err != nil {
						return func(st p.State) (interface{}, error) {
							return nil, err
						}
					}
					switch parser := pr.(type) {
					case p.P:
						return parser
					case Parsecer:
						return parser.Parser
					default:
						return func(st p.State) (interface{}, error) {
							return nil, fmt.Errorf("excpet got a parser but %v", pr)
						}
					}
				})), nil
			default:
				return nil, fmt.Errorf("excpet got a parser but %v", prs)
			}
		},
		"then": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(params[0].(Parsecer).Parser.Then(params[1].(Parsecer).Parser)), nil
		},
		"sepby1": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.SepBy1(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"sepby": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.SepBy(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"manytil": func(env Env, args ...interface{}) (Lisp, error) {
			ptype := reflect.TypeOf((*Parsecer)(nil)).Elem()
			params, err := GetArgs(env, p.UnionAll(TypeAs(ptype), TypeAs(ptype), p.EOF), args)
			if err != nil {
				return nil, err
			}
			return ParsecBox(p.ManyTil(params[0].(Parsecer).Parser, params[1].(Parsecer).Parser)), nil
		},
		"maybe": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Maybe Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Maybe(parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Manybe Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
		"skip": func(env Env, args ...interface{}) (Lisp, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Parsex Parser Skip Error: only accept one parsex parser as arg but %v", args)
			}
			param, err := Eval(env, args[0])
			if err != nil {
				return nil, err
			}
			switch parser := param.(type) {
			case Parsecer:
				return ParsecBox(p.Skip(parser.Parser)), nil
			default:
				return nil, fmt.Errorf(
					"Skip Arg Error:expect 1 parser arg but %v.",
					reflect.TypeOf(param))
			}
		},
	},
}

// NewStringState 构造一个新的基于字符串的 state
func NewStringState(data string) p.State {
	re := p.BasicStateFromText(data)
	return &re
}

// Parsecer 实现一个 parsex 封装
type Parsecer struct {
	Parser p.P
}

// Task 定义了 parsex 的求值
func (parser Parsecer) Task(env Env, args ...interface{}) (Lisp, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(
			"Parsec Parser Exprission Error: only accept one parsec state as arg but %v",
			args[0])
	}
	param, err := Eval(env, args[0])
	if err != nil {
		return nil, err
	}
	var st p.State
	var ok bool
	if st, ok = param.(p.State); !ok {
		return nil, fmt.Errorf(
			"Parsec Parser Exprission Error: only accept one parsec state as arg but %v",
			reflect.TypeOf(args[0]))
	}
	return ParsecTask{parser.Parser, st}, nil
}

// Eval 定义了其解析求值时直接返回 parser
func (parser Parsecer) Eval(env Env) (interface{}, error) {
	return parser, nil
}

// ParsecBox 定义了一个 Parsecer 的封装
func ParsecBox(parser p.P) Lisp {
	return Parsecer{parser}
}

// ParsecTask 定义了延迟执行 Parsex 的行为
type ParsecTask struct {
	Parser p.P
	State  p.State
}

// Eval 定义了 parsec task 的解析求值
func (pt ParsecTask) Eval(env Env) (interface{}, error) {
	return pt.Parser(pt.State)
}
