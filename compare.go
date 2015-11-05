package gisp

import (
	"fmt"
	"io"
	"reflect"
	tm "time"

	p "github.com/Dwarfartisan/goparsec2"
)

// FalseIfHasNil 实现一个 is nil 判断
func FalseIfHasNil(st p.State) (interface{}, error) {
	for {
		val, err := p.One(st)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("False If Nil Error: Found EOF.")
			}
			return nil, err
		}
		if val == nil {
			return false, err
		}
	}
}

// LessThanNil 实现三值 Less
func LessThanNil(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		val, _ := p.One(st)
		if x == nil || val == nil {
			return false, nil
		}
		return nil, st.Trap("expect nil value but: %v", val)
	}
}

// ListValue 检查 state 中下一个值是否是列表
func ListValue(st p.State) (interface{}, error) {
	val, err := p.One(st)
	if err == nil {
		if _, ok := val.(List); ok {
			return val, nil
		}
		return nil, fmt.Errorf("expect a list value but %v ", val)
	}
	return nil, fmt.Errorf("expect a list value but error: %v", err)
}

// LessThanList 从最近的环境中找到 < 的实现并调用其进行比较，这样用户可以自己实现特化的比较
func LessThanList(env Env) func(x interface{}) p.P {
	lessp, ok := env.Lookup("<")
	return func(x interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			if !ok {
				return nil, fmt.Errorf("Less Than List Error: opreator < not found")
			}
			y, err := ListValue(st)
			if err != nil {
				return nil, err
			}
			for _, item := range ZipLess(x.(List), y.(List)) {
				b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
				if err != nil {
					return nil, err
				}
				if b.(bool) {
					return true, nil
				}
			}
			return len(x.(List)) < len(y.(List)), nil
		}
	}
}

// LessThanListOption 允许返回 nil
func LessThanListOption(env Env) func(x interface{}) p.P {
	lessp, ok := env.Lookup("<?")
	return func(x interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			if !ok {
				return nil, fmt.Errorf("Less Than List Option Error: <? opreator not found")
			}
			y, err := ListValue(st)
			if err != nil {
				return nil, err
			}
			for _, item := range ZipLess(x.(List), y.(List)) {
				b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
				if err != nil {
					return nil, err
				}
				if b.(bool) {
					return true, nil
				}
			}
			return len(x.(List)) < len(y.(List)), nil
		}
	}
}

// TimeValue 判断 state 中下一个元素是否为 time.Time
func TimeValue(st p.State) (interface{}, error) {
	val, err := p.One(st)
	if err == nil {
		if _, ok := val.(tm.Time); ok {
			return val, nil
		}
		return nil, fmt.Errorf("expect a time value but: %v", err)
	}
	return nil, fmt.Errorf("expect a time value but error: %v", err)
}

// LessThanTime 对 Time 值进行比较
func LessThanTime(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := TimeValue(st)
		if err == nil {
			return x.(tm.Time).Before(y.(tm.Time)), nil
		}
		return nil, err
	}
}

// StringValue 判断 state 中下一个值是否为 String
func StringValue(st p.State) (interface{}, error) {
	return p.Do(func(state p.State) interface{} {
		val := p.P(p.One).Exec(state)
		if _, ok := val.(string); ok {
			return val
		}
		panic(st.Trap("expect a string but %v", val))
	})(st)
}

// LessThanInt 实现整数的比较
func LessThanInt(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := IntValue(st)
		if err == nil {
			return x.(Int) < y.(Int), nil
		}
		return nil, err
	}
}

// LessThanFloat 实现浮点数的比较
func LessThanFloat(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := FloatValue(st)
		if err == nil {
			switch val := x.(type) {
			case Float:
				return val < y.(Float), nil
			case Int:
				return Float(val) < y.(Float), nil
			default:
				return nil, st.Trap("unknown howto compoare %v < %v", x, y)
			}
		}
		return nil, err
	}
}

// LessThanNumber 实现数值的比较
func LessThanNumber(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		tran := st.Pos()
		cmp, err := LessThanInt(x)(st)
		if err == nil {
			st.Commit(tran)
			return cmp, nil
		}
		st.Rollback(tran)
		return LessThanFloat(x)(st)
	}
}

// LessThanString 实现字符串的比较
func LessThanString(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := StringValue(st)
		if err == nil {
			return x.(string) < y.(string), nil
		}
		return nil, st.Trap("expect less compare string %v and %v but error: %v",
			x, y, err)
	}
}

func lessListIn(env Env, x, y List) (interface{}, error) {
	lessp, ok := env.Lookup("<")
	if !ok {
		return nil, fmt.Errorf("Less Than List Error: < opreator not found")
	}
	for _, item := range ZipLess(x, y) {
		b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
		if err != nil {
			return nil, err
		}
		if b.(bool) {
			return true, nil
		}
	}
	return len(x) < len(y), nil
}

func lessListOptIn(env Env, x, y List) (interface{}, error) {
	lessp, ok := env.Lookup("<?")
	if !ok {
		return nil, fmt.Errorf("Less Than Option List Error: opreator <? not found")
	}
	for _, item := range ZipLess(x, y) {
		b, err := Eval(env, L(lessp, item.(List)[0], item.(List)[1]))
		if err != nil {
			return nil, err
		}
		if b.(bool) {
			return true, nil
		}
	}
	return len(x) < len(y), nil
}

func less(env Env) p.P {
	return func(st p.State) (interface{}, error) {
		l, err := p.Choice(
			p.Try(p.P(IntValue).Bind(LessThanNumber)),
			p.Try(p.P(NumberValue).Bind(LessThanFloat)),
			p.Try(p.P(StringValue).Bind(LessThanString)),
			p.Try(p.P(TimeValue).Bind(LessThanTime)),
			p.P(ListValue).Bind(LessThanList(env)),
		).Bind(func(l interface{}) p.P {
			return func(st p.State) (interface{}, error) {
				_, err := p.EOF(st)
				if err != nil {
					return nil, st.Trap("less args sign error: expect eof")
				}
				return l, nil
			}
		})(st)
		if err == nil {
			return l, nil
		}
		return nil, st.Trap("expect two lessable values compare but error %v", err)
	}
}

// return false, true or nil
func lessOption(env Env) p.P {
	return func(st p.State) (interface{}, error) {
		l, err := p.Choice(
			p.Try(p.P(IntValue).Bind(LessThanNumber)),
			p.Try(p.P(NumberValue).Bind(LessThanFloat)),
			p.Try(p.P(StringValue).Bind(LessThanString)),
			p.Try(p.P(TimeValue).Bind(LessThanTime)),
			p.Try(p.P(ListValue).Bind(LessThanListOption(env))),
			p.P(p.One).Bind(LessThanNil),
		).Bind(func(l interface{}) p.P {
			return func(st p.State) (interface{}, error) {
				_, err := p.EOF(st)
				if err != nil {
					return nil, st.Trap("less args sign error: expect eof")
				}
				return l, nil
			}
		})(st)
		if err == nil {
			return l, nil
		}
		return nil, st.Trap("expect two lessable values or nil compare but error: %v", err)
	}
}

func cmpInt(x, y Int) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpFloat(x, y Float) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpString(x, y string) Int {
	if x < y {
		return Int(1)
	}
	if y < x {
		return Int(-1)
	}
	if x == y {
		return Int(0)
	}
	return Int(0)
}

func cmpTime(x, y tm.Time) Int {
	if x.Before(y) {
		return Int(1)
	}
	if x.After(y) {
		return Int(-1)
	}
	return Int(0)
}

func cmpListIn(env Env, x, y List) (interface{}, error) {
	ret, err := lessListIn(env, x, y)
	if err != nil {
		return nil, err
	}
	if ret.(bool) {
		return -1, nil
	}
	ret, err = lessListIn(env, y, x)
	if err != nil {
		return nil, err
	}
	if ret.(bool) {
		return 1, nil
	}
	if reflect.DeepEqual(x, y) {
		return 0, nil
	}
	return nil, fmt.Errorf("Compare Error: Unknown howto copmare %v and %v", x, y)
}

// CmpInt 实现两个整数的三向比较
func CmpInt(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := IntValue(st)
		if err == nil {
			return cmpInt(x.(Int), y.(Int)), nil
		}
		return nil, err
	}
}

// CmpFloat 实现两个浮点数的三向比较
func CmpFloat(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := FloatValue(st)
		if err == nil {
			switch val := x.(type) {
			case Float:
				return cmpFloat(val, y.(Float)), nil
			case Int:
				return cmpFloat(Float(val), y.(Float)), nil
			default:
				return nil, st.Trap("unknown howto compoare %v < %v", x, y)
			}
		}
		return nil, err
	}
}

// CmpNumber 实现两个数值的三向比较
func CmpNumber(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		tran := st.Begin()
		cmp, err := CmpInt(x)(st)
		if err == nil {
			st.Commit(tran)
			return cmp, nil
		}
		st.Rollback(tran)
		return CmpFloat(x)(st)
	}
}

// CmpString 实现两个字符串的三向比较
func CmpString(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := StringValue(st)
		if err == nil {
			return cmpString(x.(string), y.(string)), nil
		}
		return nil, st.Trap("expect less compare string %v and %v but error: %v",
			x, y, err)
	}
}

// CmpTime 实现两个Time的三向比较
func CmpTime(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := TimeValue(st)
		if err == nil {
			return cmpTime(x.(tm.Time), y.(tm.Time)), nil
		}
		return nil, fmt.Errorf("expect less compare string %v and %v but error: %v",
			x, y, err)
	}
}

func compare(st p.State) (interface{}, error) {
	l, err := p.Choice(
		p.Try(IntValue).Bind(LessThanNumber),
		p.Try(NumberValue).Bind(LessThanFloat),
		p.Try(StringValue).Bind(LessThanString),
		p.P(TimeValue).Bind(LessThanTime),
	).Bind(func(l interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			_, err := p.EOF(st)
			if err != nil {
				return nil, fmt.Errorf("less args sign error: expect eof")
			}
			return l, nil
		}
	})(st)
	if err == nil {
		return l, nil
	}
	return nil, fmt.Errorf("expect two lessable values compare but error %v", err)
}

func equals(st p.State) (interface{}, error) {
	return p.P(p.One).Bind(eqs)(st)
}
func eqs(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := st.Next()
		if err != nil {
			if e, ok := err.(p.Error); ok && e.Message == "eof" {
				return true, nil
			}
			return nil, err
		}
		if reflect.DeepEqual(x, y) {
			return eqs(x)(st)
		}
		return false, nil
	}
}

func equalsOption(st p.State) (interface{}, error) {
	return p.P(p.One).Bind(eqsOption)(st)
}

func eqsOption(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := st.Next()
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return true, nil
			}
			return nil, err
		}
		if x == nil || y == nil {
			return false, nil
		}
		if reflect.DeepEqual(x, y) {
			return eqsOption(x)(st)
		}
		return false, nil
	}
}

func notEquals(st p.State) (interface{}, error) {
	return p.P(p.One).Bind(neqs)(st)
}

func neqs(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := st.Next()
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return false, nil
			}
			return nil, err
		}
		if x == nil || y == nil {
			return false, nil
		}
		if !reflect.DeepEqual(x, y) {
			return neqs(x)(st)
		}
		return false, nil
	}
}

// not equals function, NotEqual or !=, if anyone is nil, return false
func neqsOption(st p.State) (interface{}, error) {
	x, err := st.Next()
	if err != nil {
		return nil, err
	}
	if x == nil {
		return false, nil
	}
	for {
		y, err := st.Next()
		if err != nil {
			if reflect.DeepEqual(err, io.EOF) {
				return false, nil
			}
			return nil, err
		}
		if y == nil {
			return false, nil
		}
		if !reflect.DeepEqual(x, y) {
			return true, nil
		}
	}
}

// String2Values 将两个 StringValue 串为 List
var String2Values = p.P(StringValue).Bind(func(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := StringValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

//TimeValue 将两个 Time 值串为 List
var Time2Values = p.P(TimeValue).Bind(func(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := TimeValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

//ListValue 将两个 Time 值串为 List
var List2Values = p.P(ListValue).Bind(func(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		y, err := ListValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})
