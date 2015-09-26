package gisp

import (
	"reflect"

	p "github.com/Dwarfartisan/goparsec2"
)

//Type 对象定义了一个可空的反射类型，用于 Lisp 对象定义
type Type struct {
	reflect.Type
	option bool
}

func (typ Type) String() string {
	str := typ.Type.String()
	if typ.option {
		return str + "?"
	}
	return str
}

// Option 指示类型是否是可空的
func (typ Type) Option() bool {
	return typ.option
}

func stop(st p.State) (interface{}, error) {
	pos := st.Pos()
	defer st.SeekTo(pos)
	r, err := p.Choice(
		p.Try(p.Space),
		p.Try(p.Newline),
		p.Try(p.RuneOf(":.()[]{}?")),
		p.Try(p.EOF),
	)(st)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func stopWord(x interface{}) p.Parsec {
	return p.M(stop).Then(p.Return(x))
}

func typeName(word string) p.Parsec {
	return p.Str(word).Bind(stopWord)
}

var anyType = p.Many1(p.Choice(p.Try(p.Digit), p.Letter)).Bind(stopWord).Bind(p.ReturnString)

// SliceTypeParserExt 定义了带环境的序列类型解析逻辑
func SliceTypeParserExt(env Env) p.Parsec {
	return func(st p.State) (interface{}, error) {
		t, err := p.Str("[]").Then(ExtTypeParser(env))(st)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(t.(Type).Type), nil
	}
}

// MapTypeParserExt  定义了带环境的映射类型解析逻辑
func MapTypeParserExt(env Env) p.Parsec {
	return func(st p.State) (interface{}, error) {
		key, err := p.Between(p.Str("map["), p.Chr(']'), ExtTypeParser(env))(st)
		if err != nil {
			return nil, err
		}
		value, err := ExtTypeParser(env)(st)
		if err != nil {
			return nil, err
		}
		return reflect.MapOf(key.(Type).Type, value.(Type).Type), nil
	}
}

// MapTypeParser 定义了序列类型解析逻辑
func MapTypeParser(st p.State) (interface{}, error) {
	key, err := p.Between(p.Str("map["), p.Chr(']'), TypeParser)(st)
	if err != nil {
		return nil, err
	}
	value, err := TypeParser(st)
	if err != nil {
		return nil, err
	}
	return reflect.MapOf(key.(Type).Type, value.(Type).Type), nil
}

// ExtTypeParser 定义了带环境的类型解释器
func ExtTypeParser(env Env) p.Parsec {
	return func(st p.State) (interface{}, error) {
		_, err := p.Str("::")(st)
		if err != nil {
			return nil, err
		}
		buildin := p.Choice(
			p.Try(typeName("bool").Then(p.Return(BOOL))),
			p.Try(typeName("float").Then(p.Return(FLOAT))),
			p.Try(typeName("int").Then(p.Return(INT))),
			p.Try(typeName("string").Then(p.Return(STRING))),
			p.Try(typeName("time").Then(p.Return(TIME))),
			p.Try(typeName("duration").Then(p.Return(DURATION))),
			p.Try(typeName("any").Then(p.Return(ANY))),
			p.Try(typeName("atom").Then(p.Return(ATOM))),
			p.Try(p.Str("list").Then(p.Return(LIST))),
			p.Try(typeName("quote").Then(p.Return(QUOTE))),
			p.Try(p.Str("dict").Then(p.Return(DICT))),
			p.Try(MapTypeParserExt(env)),
		)
		ext := func(st p.State) (interface{}, error) {
			n, err := anyType(st)
			if err != nil {
				return nil, err
			}
			t, ok := env.Lookup(n.(string))
			if !ok {
				return nil, st.Trap("type %v not found", n)
			}
			if typ, ok := t.(reflect.Type); ok {
				return typ, nil
			}
			return nil, st.Trap("var %v is't a type. It is %v", n, reflect.TypeOf(t))
		}
		t, err := p.Choice(buildin, ext)(st)
		if err != nil {
			return nil, err
		}
		_, err = p.Try(p.Chr('?'))(st)
		option := err == nil
		return Type{t.(reflect.Type), option}, nil
	}
}

// TypeParser 定义了一个基本的类型解释器
func TypeParser(st p.State) (interface{}, error) {
	t, err := p.Str("::").Then(
		p.Choice(
			p.Try(p.Str("bool").Then(p.Return(BOOL))),
			p.Try(p.Str("float").Then(p.Return(FLOAT))),
			p.Try(p.Str("int").Then(p.Return(INT))),
			p.Try(p.Str("string").Then(p.Return(STRING))),
			p.Try(p.Str("time").Then(p.Return(TIME))),
			p.Try(p.Str("duration").Then(p.Return(DURATION))),
			p.Try(p.Str("any").Then(p.Return(ANY))),
			p.Try(p.Str("atom").Then(p.Return(ATOM))),
			p.Try(p.Str("list").Then(p.Return(LIST))),
			p.Try(p.Str("quote").Then(p.Return(QUOTE))),
			p.Try(p.Str("dict").Then(p.Return(DICT))),
			MapTypeParser,
		))(st)
	if err != nil {
		return nil, err
	}
	_, err = p.Try(p.Chr('?'))(st)
	option := err == nil
	return Type{t.(reflect.Type), option}, nil
}
