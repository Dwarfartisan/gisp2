package gisp

import (
	"fmt"

	p "github.com/Dwarfartisan/goparsec2"
)

// Bool 是内置的 bool 类型的封装
type Bool bool

// BoolParser 解析 bool
var BoolParser = p.Choice(p.Str("true"), p.Str("false")).Bind(func(input interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		switch input.(string) {
		case "true":
			return Bool(true), nil
		case "false":
			return Bool(false), nil
		default:
			return nil, fmt.Errorf("Unexpect bool token %v", input)
		}
	}
})

// NilParser 解析 nil
var NilParser = p.Str("nil").Then(p.Return(nil))

// Nil 类型定义空值行为
type Nil struct {
}

// Eval 方法实现 Lisp 值的 Eval 。 Nil 的 Eval 总是返回空值
func (n Nil) Eval(env Env) (interface{}, error) {
	return nil, nil
}
