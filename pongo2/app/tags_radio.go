package pongo2

import (
	"fmt"
	"html"

	"github.com/revel/revel"
	p2 "github.com/runner-mei/pongo2"
)

type tagRadioNode struct {
	field string
	value p2.IEvaluator
}

func (node *tagRadioNode) Execute(ctx *p2.ExecutionContext, writer p2.TemplateWriter) *p2.Error {
	fieldObj := ctx.Public[node.field]
	if nil == fieldObj {
		return ctx.Error("field '"+node.field+"' is missing.", nil)
	}
	field, _ := fieldObj.(*revel.Field)
	if nil == field {
		return ctx.Error(fmt.Sprintf("field '"+node.field+"' isn't Field - %T.", fieldObj), nil)
	}
	val, err := node.value.Evaluate(ctx)
	if err != nil {
		return err
	}
	val_str := val.String()

	checked := ""
	if field.Flash() == val_str {
		checked = " checked"
	}
	fmt.Fprintf(writer, `<input type="radio" name="%s" value="%s"%s>`,
		html.EscapeString(field.Name), html.EscapeString(val_str), checked)
	return nil
}

// tagURLForParser implements a {% urlfor %} tag.
//
// urlfor takes one argument for the controller, as well as any number of key/value pairs for additional URL data.
// Example: {% urlfor "UserController.View" ":slug" "oal" %}
func tagRadioParser(doc *p2.Parser, start *p2.Token, arguments *p2.Parser) (p2.INodeTag, *p2.Error) {
	var field string
	typeToken := arguments.MatchType(p2.TokenIdentifier)
	if typeToken != nil {
		field = typeToken.Val
	} else if sToken := arguments.MatchType(p2.TokenString); nil != sToken {
		field = sToken.Val
	} else {
		return nil, arguments.Error("Expected an identifier or string.", nil)
	}

	expr, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &tagRadioNode{field: field,
		value: expr}, nil
}

func init() {
	p2.RegisterTag("radio", tagRadioParser)
}
