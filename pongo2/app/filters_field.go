package pongo2

import (
	"github.com/revel/revel"
	p2 "github.com/runner-mei/pongo2"
)

func init() {
	p2.RegisterFilter("field", func(ctx *p2.ExecutionContext, in *p2.Value, param *p2.Value) (out *p2.Value, err *p2.Error) {
		if nil == in.Interface() || in.String() == "" {
			return nil, ctx.Error("field argument must is string", nil)
		}
		return p2.AsValue(revel.NewField(in.String(), ctx.Public)), nil
	})
}
