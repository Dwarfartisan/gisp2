package gisp

import (
	//px "github.com/Dwarfartisan/goparsec/parsex"
	"fmt"
)

// Utils 提供了一个常用函数工具包
var Utils = Toolkit{
	Meta: map[string]interface{}{
		"name":     "utils",
		"category": "package",
	},
	Content: map[string]interface{}{
		"errorf": func(env Env, args ...interface{}) (Tasker, error) {
			if len(args) < 1 {
				return nil, fmt.Errorf("Errorf Empty Arg Error:expect args has 1 arg a last.")
			}
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			if tmpl, ok := params[0].(string); ok {
				return func(env Env) (interface{}, error) {
					return nil, fmt.Errorf(tmpl, params[1:]...)
				}, nil
			}
			return nil, fmt.Errorf("Errorf Arg Error:expect first arg is a string but %v.", args[0])
		},
		"error": func(env Env, args ...interface{}) (Tasker, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("Error Arg Error:expect args has 1 arg.")
			}
			params, err := Evals(env, args...)
			if err != nil {
				return nil, err
			}
			return func(env Env) (interface{}, error) {
				return nil, params[0].(error)
			}, nil
		},
		"printf": printf,
		"ginq": LispExpr(func(env Env, args ...interface{}) (Lisp, error) {
			return Q(NewGinq(args...)), nil
		}),
	},
}

func printf(env Env, args ...interface{}) (Tasker, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("Printf Empty Arg Error:expect args has 1 arg a last.")
	}
	params, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	if tmpl, ok := params[0].(string); ok {
		return func(env Env) (interface{}, error) {
			return fmt.Printf(tmpl, params[1:]...)
		}, nil
	}
	return nil, fmt.Errorf("Printf Arg Error:expect first arg is a string but %v.", args[0])
}
