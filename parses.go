package gisp

import (
	"fmt"
	"reflect"
	"strconv"

	p "github.com/Dwarfartisan/goparsec2"
)

// Ext 扩展表示扩展环境

// Skip 忽略匹配指定算子的内容
var Skip = p.Skip(p.Space)

// IntParser 解析整数
func IntParser(st p.State) (interface{}, error) {
	i, err := p.Int(st)
	if err == nil {
		val, err := strconv.Atoi(i.(string))
		if err == nil {
			return Int(val), nil
		}
		return nil, err
	}
	return nil, err

}

// 用于string
var EscapeChars = p.Do(func(st p.State) interface{} {
	p.Chr('\\').Exec(st)
	r := p.RuneOf("nrt\"\\").Exec(st)
	ru := r.(rune)
	switch ru {
	case 'r':
		return '\r'
	case 'n':
		return '\n'
	case '"':
		return '"'
	case '\\':
		return '\\'
	case 't':
		return '\t'
	default:
		panic(st.Trap("Unknown escape sequence \\%c", r))
	}
})

//用于rune
var EscapeCharr = p.Do(func(st p.State) interface{} {
	p.Chr('\\').Exec(st)
	r := p.RuneOf("nrt'\\").Exec(st)
	ru := r.(rune)
	switch ru {
	case 'r':
		return '\r'
	case 'n':
		return '\n'
	case '\'':
		return '\''
	case '\\':
		return '\\'
	case 't':
		return '\t'
	default:
		panic(st.Trap("Unknown escape sequence \\%c", r))
	}
})

// RuneParser 实现 rune 的解析
var RuneParser = p.Do(func(state p.State) interface{} {
	p.Chr('\'').Exec(state)
	c := p.Choice(p.Try(EscapeCharr), p.NChr('\'')).Exec(state)
	p.Chr('\'').Exec(state)
	return Rune(c.(rune))
})

// StringParser 实现字符串解析
var StringParser = p.Between(p.Chr('"'), p.Chr('"'),
	p.Many(p.Choice(p.Try(EscapeChars), p.NChr('"')))).Bind(p.ReturnString)

func bodyParser(st p.State) (interface{}, error) {
	value, err := p.SepBy(ValueParser(), Skip)(st)
	return value, err
}

func bodyParserExt(env Env) p.Parsec {
	return p.Many(ValueParserExt(env).Over(Skip))
}

// ListParser 实现列表解析器
func ListParser() p.Parsec {
	return func(st p.State) (interface{}, error) {
		left := p.Chr('(').Then(Skip)
		right := Skip.Then(p.Chr(')'))
		empty := p.Between(left, right, Skip)
		list, err := p.Between(left, right, bodyParser)(st)
		if err == nil {
			switch l := list.(type) {
			case List:
				return L(l), nil
			case []interface{}:
				return list.([]interface{}), nil
			default:
				return nil, fmt.Errorf("List Parser Error: %v type is unexpected: %v", list, reflect.TypeOf(list))
			}
		} else {
			_, e := empty(st)
			if e == nil {
				return List{}, nil
			}
			return nil, err
		}
	}
}

// ListParserExt 实现带扩展的列表解析器
func ListParserExt(env Env) p.Parsec {
	left := p.Chr('(').Then(Skip)
	right := Skip.Then(p.Chr(')'))
	empty := left.Then(right)
	return func(st p.State) (interface{}, error) {
		list, err := p.Try(p.Between(left, right, bodyParserExt(env)))(st)
		if err == nil {
			switch l := list.(type) {
			case List:
				return L(l), nil
			case []interface{}:
				return List(l), nil
			default:
				return nil, fmt.Errorf("List Parser(ext) Error: %v type is unexpected: %v", list, reflect.TypeOf(list))
			}
		} else {
			_, e := empty(st)
			if e == nil {
				return List{}, nil
			}
			return nil, err
		}
	}
}

// QuoteParser 实现 Quote 语法的解析
func QuoteParser(st p.State) (interface{}, error) {
	lisp, err := p.Chr('\'').Then(
		p.Choice(
			p.Try(p.M(AtomParser).Bind(SuffixParser)),
			ListParser().Bind(SuffixParser),
		))(st)
	if err == nil {
		return Quote{lisp}, nil
	}
	return nil, err
}

// QuoteParserExt 实现带扩展的 Quote 语法的解析
func QuoteParserExt(env Env) p.Parsec {
	return func(st p.State) (interface{}, error) {
		lisp, err := p.Chr('\'').Then(p.Choice(
			p.Try(AtomParserExt(env).Bind(SuffixParser)),
			ListParserExt(env).Bind(SuffixParser),
		))(st)
		if err == nil {
			return Quote{lisp}, nil
		}
		return nil, err
	}
}

// ValueParser 实现简单的值解释器
func ValueParser() p.Parsec {
	return func(state p.State) (interface{}, error) {
		value, err := p.Choice(p.Try(StringParser),
			p.Try(FloatParser),
			p.Try(IntParser),
			p.Try(RuneParser),
			p.Try(StringParser),
			p.Try(BoolParser),
			p.Try(NilParser),
			p.Try(p.M(AtomParser).Bind(SuffixParser)),
			p.Try(p.M(ListParser()).Bind(SuffixParser)),
			p.Try(DotExprParser),
			QuoteParser,
		)(state)
		return value, err
	}
}

// ValueParserExt 表示带扩展的值解释器
func ValueParserExt(env Env) p.Parsec {
	return func(st p.State) (interface{}, error) {
		value, err := p.Choice(p.Try(StringParser),
			p.Try(FloatParser),
			p.Try(IntParser),
			p.Try(RuneParser),
			p.Try(StringParser),
			p.Try(BoolParser),
			p.Try(NilParser),
			p.Try(AtomParserExt(env).Bind(SuffixParserExt(env))),
			p.Try(ListParserExt(env).Bind(SuffixParserExt(env))),
			p.Try(DotExprParser),
			p.Try(BracketExprParserExt(env)),
			QuoteParserExt(env),
		)(st)
		return value, err
	}
}
