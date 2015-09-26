package gisp

import (
	"strconv"

	p "github.com/Dwarfartisan/goparsec2"
)

// Float 是 gisp 系统的浮点数实现
type Float float64

// FloatParser 解析浮点数
func FloatParser(state p.State) (interface{}, error) {
	return p.Do(func(st p.State) interface{} {
		f := p.Try(p.Float).Exec(st)
		val, err := strconv.ParseFloat(f.(string), 64)
		if err == nil {
			return Float(val)
		}
		panic(st.Trap("%v", err))
	})(state)
}
