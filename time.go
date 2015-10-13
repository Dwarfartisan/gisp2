package gisp

import (
	tm "time"

	p "github.com/Dwarfartisan/goparsec2"
)

// Time 包引入了go的time包功能
var Time = Toolkit{
	Meta: map[string]interface{}{
		"category": "toolkit",
		"name":     "time",
	},
	Content: map[string]interface{}{
		"now": SimpleBox{
			SignChecker(p.EOF),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.Now(), nil
				}
			}},
		"parseDuration": SimpleBox{
			SignChecker(p.P(StringValue).Then(p.EOF)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.ParseDuration(args[0].(string))
				}
			}},
		"parseTime": SimpleBox{
			SignChecker(p.P(StringValue).Then(StringValue).Then(p.EOF)),
			func(args ...interface{}) Tasker {
				return func(env Env) (interface{}, error) {
					return tm.Parse(args[0].(string), args[1].(string))
				}
			}},
	},
}
