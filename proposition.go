package gisp

import (
	"fmt"

	p "github.com/Dwarfartisan/goparsec2"
)

// Propositions 给出了一组常用的操作
var Propositions = Toolkit{
	Meta: map[string]interface{}{
		"name":     "propositions",
		"category": "package",
	},
	Content: map[string]interface{}{
		"lambda": BoxExpr(LambdaExpr),
		"let":    BoxExpr(LetExpr),
		"+":      EvalExpr(ParsecExpr(addx)),
		"add":    EvalExpr(ParsecExpr(addx)),
		"-":      EvalExpr(ParsecExpr(subx)),
		"sub":    EvalExpr(ParsecExpr(subx)),
		"*":      EvalExpr(ParsecExpr(mulx)),
		"mul":    EvalExpr(ParsecExpr(mulx)),
		"/":      EvalExpr(ParsecExpr(divx)),
		"div":    EvalExpr(ParsecExpr(divx)),
		"cmp":    EvalExpr(cmpExpr),
		"less":   EvalExpr(lessExpr),
		"<":      EvalExpr(lessExpr),
		"<?":     EvalExpr(lsoExpr),
		"<=":     EvalExpr(leExpr),
		"<=?":    EvalExpr(leoExpr),
		">":      EvalExpr(greatExpr),
		">?":     EvalExpr(gtoExpr),
		">=":     EvalExpr(geExpr),
		">=?":    EvalExpr(geoExpr),
		"==":     EvalExpr(eqsExpr),
		"==?":    EvalExpr(eqsoExpr),
		"!=":     EvalExpr(neqsExpr),
		"!=?":    EvalExpr(neqsoExpr),
	},
}

// ParsecExpr 是 parsec 算子的解析表达式
func ParsecExpr(pxExpr p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		st := p.NewBasicState(data)
		ret, err := pxExpr(&st)
		if err != nil {
			return nil, err
		}
		return Q(ret), nil
	}
}

// ExtExpr 带扩展环境
func ExtExpr(extExpr func(env Env) p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		st := p.NewBasicState(data)
		ret, err := extExpr(env)(&st)
		if err != nil {
			return nil, err
		}
		return Q(ret), nil
	}
}

// NotParsec 是 not 运算符
func NotParsec(pxExpr p.P) p.P {
	return func(st p.State) (interface{}, error) {
		b, err := pxExpr(st)
		if err != nil {
			return nil, err
		}
		if boolean, ok := b.(bool); ok {
			return !boolean, nil
		}
		return nil, st.Trap("Unknow howto not %v", b)
	}
}

// ParsecReverseExpr 是倒排运算
func ParsecReverseExpr(pxExpr p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		l := len(data)
		last := l - 1
		datax := make([]interface{}, l)
		for idx, item := range data {
			datax[last-idx] = item
		}
		st := p.NewBasicState(data)
		x, err := pxExpr(&st)
		if err != nil {
			return nil, err
		}
		return Q(x), nil
	}
}

// NotExpr 定义了 not 表达式
func NotExpr(expr LispExpr) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		element, err := expr(env, args...)
		if err != nil {
			return nil, err
		}
		ret, err := element.Eval(env)
		if err != nil {
			return nil, err
		}
		var b bool
		if b, ok := ret.(bool); ok {
			return Q(!b), nil
		}
		return nil, fmt.Errorf("Unknow howto not %v", b)
	}
}

// OrExpr 是  or 表达式
func OrExpr(x, y p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		data, err := Evals(env, args...)
		if err != nil {
			return nil, err
		}
		st := p.NewBasicState(data)
		rex, err := x(&st)
		if err != nil {
			fmt.Println("Trace x parsec")
			return nil, err
		}
		if b, ok := rex.(bool); ok {
			if b {
				return Q(true), nil
			}
			st.SeekTo(0)
			rex, err = y(&st)
			if err != nil {
				fmt.Println("Trace y parsec")
				return nil, err
			}
			return Q(rex), nil
		}
		return nil, fmt.Errorf("Unknow howto combine %v or %v for %v", x, y, data)
	}
}

// OrExtExpr 定了带环境扩展的 or 表达式
func OrExtExpr(x, y func(Env) p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return OrExpr(x(env), y(env))(env, args...)
	}
}

// OrExtRExpr 定了带环境扩展的 or 逆向表达式
func OrExtRExpr(x p.P, y func(Env) p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return OrExpr(x, y(env))(env, args...)
	}
}

// ExtReverseExpr 定了带环境扩展的倒排表达式
func ExtReverseExpr(expr func(Env) p.P) LispExpr {
	return func(env Env, args ...interface{}) (Lisp, error) {
		return ParsecReverseExpr(expr(env))(env, args...)
	}
}

var addExpr = ParsecExpr(addx)
var subExpr = ParsecExpr(subx)
var mulExpr = ParsecExpr(mulx)
var divExpr = ParsecExpr(divx)
var lessExpr = ExtExpr(less)
var lsoExpr = ExtExpr(lessOption)
var leExpr = OrExtRExpr(equals, less)
var leoExpr = OrExtRExpr(equalsOption, lessOption)
var cmpExpr = ParsecExpr(compare)
var greatExpr = ExtReverseExpr(less)
var gtoExpr = ExtReverseExpr(lessOption)
var geExpr = OrExtRExpr(equals, less)
var geoExpr = func(env Env, args ...interface{}) (Lisp, error) {
	st := p.NewBasicState(args)
	ret, err := p.Choice(p.Try(NotParsec(less(env))), FalseIfHasNil)(&st)
	if err != nil {
		return nil, err
	}
	return Q(ret), nil
}
var eqsExpr = ParsecExpr(equals)
var eqsoExpr = ParsecExpr(equalsOption)
var neqsExpr = NotExpr(eqsExpr)
var neqsoExpr = ParsecExpr(neqsOption)
