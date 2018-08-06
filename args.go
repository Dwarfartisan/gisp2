package gisp

import (
	"fmt"
	"reflect"

	p "github.com/Dwarfartisan/goparsec2"
)

func typeis(x Atom) func(int, interface{}) (interface{}, error) {
	return func(pos int, data interface{}) (interface{}, error) {
		if data == nil {
			if x.Type.Option() {
				return data, nil
			}
			return nil, fmt.Errorf("%v's type not match %v", data, x.Type)
		}
		if reflect.DeepEqual(x.Type.Type, ANY) {
			return data, nil
		}
		if reflect.DeepEqual(x.Type.Type, reflect.TypeOf(data)) {
			return data, nil
		}
		return data, TypeSignError{x.Type, data}
	}
}

// TypeAs 函数根据反射对 Gisp 数据进行类型判断
func TypeAs(typ reflect.Type) p.P {
	return func(st p.State) (interface{}, error) {
		obj, err := st.Next()
		if err != nil {
			return nil, err
		}
		otype := reflect.TypeOf(obj)
		if otype == typ {
			return obj, nil
		}
		return nil, fmt.Errorf("Args Type Sign Check: excpet %v but %v is %v",
			typ, obj, otype)
	}
}

// argParser 构造一个 parsec 解析器，判断输入数据是否与给定类型一致，如果判断成功，构造对应的
// Var。
func argParser(atom Atom) p.P {
	one := func(st p.State) (interface{}, error) {
		var err error
		if data, err := st.Next(); err == nil {
			if _, err := typeis(atom)(st.Pos(), data); err == nil {
				slot := VarSlot(atom.Type)
				slot.Set(data)
				return slot, nil
			}
		}
		return nil, err
	}
	if atom.Name == "..." {
		return p.Many(one)
	}
	return one
}

// argRing 组成参数解析链的的后续逻辑，供 parsex.Binds 调用
func argRing(atom Atom) func(interface{}) p.P {
	return func(x interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			ring, err := argParser(atom)(st)
			if err == nil {
				return append(x.([]Var), ring.([]Var)...), nil
			}
			return nil, err
		}
	}
}

// GetArgs 方法为将传入的 args 的 gisp 值从指定环境中解析出来，然后传入 parser 。
func GetArgs(env Env, parser p.P, args []interface{}) ([]interface{}, error) {
	ret, err := Evals(env, args...)
	if err != nil {
		return nil, err
	}
	st := p.NewBasicState(ret)
	_, err = parser(&st)
	if err != nil {
		return nil, fmt.Errorf("Args Type Sign Check got error:%v", err)
	}
	return ret, nil
}
