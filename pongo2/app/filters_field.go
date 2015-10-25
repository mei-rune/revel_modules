package pongo2

import (
	p2 "github.com/flosch/pongo2"
	"github.com/revel/revel"
)

func init() {
	p2.RegisterFilter("field", func(ctx *p2.ExecutionContext, in *p2.Value, param *p2.Value) (out *p2.Value, err *p2.Error) {
		if nil == in.Interface() || in.String() == "" {
			return nil, ctx.Error("field argument must is string", nil)
		}
		return p2.AsValue(revel.NewField(in.String(), ctx.Public)), nil
	})
}
