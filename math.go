package gisp

import (
	"fmt"
	"reflect"

	p "github.com/Dwarfartisan/goparsec2"
)

// TypeMatchError 定了类型匹配错误
type TypeMatchError struct {
	Value interface{}
	Type  reflect.Type
}

func (err TypeMatchError) Error() string {
	return fmt.Sprintf("%v not match type %v", err.Value, err.Type)
}

// NotIntError 定了预期为整型但是校验失败的错误
type NotIntError struct {
	Value interface{}
}

func (err NotIntError) Error() string {
	return fmt.Sprintf("%v is't a valid Int", err.Value)
}

// NotFloatError 定义了预期为浮点数但是校验失败的错误
type NotFloatError struct {
	Value interface{}
}

func (err NotFloatError) Error() string {
	return fmt.Sprintf("%v is't a valid Float", err.Value)
}

// NotNumberError 定义了预期为数值但是校验失败的错误
type NotNumberError struct {
	Value interface{}
}

func (err NotNumberError) Error() string {
	return fmt.Sprintf("%v is't a valid Number", err.Value)
}

// IntValue 将所有整型处理为 Int ，其它类型不接受
func IntValue(st p.State) (interface{}, error) {
	v, err := st.Next()
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case int:
		return Int(val), nil
	case int8:
		return Int(val), nil
	case int16:
		return Int(val), nil
	case int32:
		return Int(val), nil
	case int64:
		return Int(val), nil
	case Int:
		return val, nil
	default:
		return nil, NotIntError{v}
	}
}

// FloatValue 将所有浮点型处理为 Float ，其它类型不接受
func FloatValue(st p.State) (interface{}, error) {
	v, err := st.Next()
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case float32:
		return Float(val), nil
	case float64:
		return Float(val), nil
	case Float:
		return val, nil
	default:
		return nil, NotFloatError{v}
	}
}

// NumberValue 将所有整型和浮点型处理为 Float ，其它类型不接受
func NumberValue(st p.State) (interface{}, error) {
	v, err := st.Next()
	if err != nil {
		return nil, err
	}
	switch val := v.(type) {
	case int:
		return Float(val), nil
	case int8:
		return Float(val), nil
	case int16:
		return Float(val), nil
	case int32:
		return Float(val), nil
	case int64:
		return Float(val), nil
	case Int:
		return Float(val), nil
	case float32:
		return Float(val), nil
	case float64:
		return Float(val), nil
	case Float:
		return val, nil
	default:
		return nil, NotNumberError{v}
	}
}

// Int2Values 获取多个 int
var Int2Values = p.M(IntValue).Bind(func(x interface{}) p.Parsec {
	return func(st p.State) (interface{}, error) {
		y, err := IntValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

// Num2Values 获取多个 int
var Num2Values = p.M(NumberValue).Bind(func(x interface{}) p.Parsec {
	return func(st p.State) (interface{}, error) {
		y, err := NumberValue(st)
		if err != nil {
			return nil, err
		}
		return []interface{}{x, y}, nil
	}
})

func xEOF(x interface{}) p.Parsec {
	return p.M(p.EOF).Then(p.Return(x))
}

// addx 实现一个parsex累加解析器，精度向上适配。我一直觉得应该有一个简单的高效版本，不需要回溯的
// 但是目前还没有找到。
func addx(st p.State) (interface{}, error) {
	ints, err := p.Try(p.ManyTil(IntValue, p.EOF))(st)
	if err == nil {
		root := Int(0)
		for _, x := range ints.([]interface{}) {
			root += x.(Int)
		}
		return root, nil
	}
	numbers, err := p.ManyTil(NumberValue, p.EOF)(st)
	if err == nil {
		root := Float(0)
		for _, x := range numbers.([]interface{}) {
			root += x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: FLOATMUST, Value: nerr.Value}
	}
	return nil, err
}

func addInts(ints ...interface{}) (interface{}, error) {
	root := ints[0].(Int)
	for _, i := range ints[1:] {
		root += i.(Int)
	}
	return root, nil
}

func addFloats(floats ...interface{}) (interface{}, error) {
	root := floats[0].(Float)
	for _, f := range floats[1:] {
		root += f.(Float)
	}
	return root, nil
}

// subx 实现一个左折叠的 parsex 连减解析器，精度向上适配。
func subx(st p.State) (interface{}, error) {
	data, err := p.Try(p.ManyTil(IntValue, p.EOF))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root -= x.(Int)
		}
		return root, nil
	}
	data, err = p.ManyTil(NumberValue, p.EOF)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root -= x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}

// mulx 实现一个 parsec 累乘解析器，精度向上适配。
func mulx(st p.State) (interface{}, error) {
	data, err := p.Try(p.ManyTil(IntValue, p.EOF))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root *= x.(Int)
		}
		return root, nil
	}
	data, err = p.ManyTil(NumberValue, p.EOF)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root *= x.(Float)
		}
		return root, nil
	}
	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}

// divx 实现一个左折叠的 parsec 连除解析器，精度向上适配。
func divx(st p.State) (interface{}, error) {
	data, err := p.Try(p.ManyTil(IntValue, p.EOF))(st)
	if err == nil {
		ints := data.([]interface{})
		root := ints[0].(Int)
		for _, x := range ints[1:] {
			root /= x.(Int)
		}
		return root, nil
	}
	data, err = p.ManyTil(NumberValue, p.EOF)(st)
	if err == nil {
		numbers := data.([]interface{})
		root := numbers[0].(Float)
		for _, x := range numbers[1:] {
			root /= x.(Float)
		}
		return root, nil
	}

	if nerr, ok := err.(NotNumberError); ok {
		return nil, TypeSignError{Type: Type{FLOAT, false}, Value: nerr.Value}
	}
	return nil, err
}
