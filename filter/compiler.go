package filter

import (
	"fmt"
	"reflect"
)

type Visitor interface {
	VisitBool(*node)
	VisitNumber(*node)
	VisitString(*node)
	VisitField(*node)
	VisitMul(lhs, rhs *node)
	VisitDiv(lhs, rhs *node)
	VisitMod(lhs, rhs *node)
	VisitAdd(lhs, rhs *node)
	VisitSub(lhs, rhs *node)
	VisitEq(lhs, rhs *node)
	VisitNe(lhs, rhs *node)
	VisitLt(lhs, rhs *node)
	VisitLe(lhs, rhs *node)
	VisitGt(lhs, rhs *node)
	VisitGe(lhs, rhs *node)
	VisitHas(field, value *node)
	VisitNot(*node)
	VisitAnd(children ...*node)
	VisitOr(children ...*node)
}

func compile(root *node, element reflect.Type) (*program, error) {
	if element.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid element kind: %v", element.Kind())
	}

	c := &compiler{
		program: &program{
			instructions: nil,
			element:      element,
		},
	}

	root.visit(c)

	return c.program, nil
}

type compiler struct {
	program *program
}

func (c *compiler) VisitBool(n *node) {
	c.program.append(instr{
		Type: instPushBool,
		Arg: argument{
			asBool: n.asBool,
		},
	})
}

func (c *compiler) VisitNumber(n *node) {
	c.program.append(instr{
		Type: instrPushNumber,
		Arg: argument{
			asInt:   n.asInt,
			asFloat: n.asFloat,
		},
	})
}

func (c *compiler) VisitString(n *node) {
	c.program.append(instr{
		Type: instrPushString,
		Arg: argument{
			asString: n.asString,
		},
	})
}

func (c *compiler) VisitField(n *node) {
	var index int64

	// TODO(daniel): better field lookup (json tag?)
	// TODO(daniel): do we need to store the type?
	for i := 0; i < c.program.element.NumField(); i++ {
		field := c.program.element.Field(i)
		if field.Name == n.asField {
			index = int64(i)
			break
		}
	}

	c.program.append(instr{
		Type: instrPushField,
		Arg: argument{
			asString: n.asField,
			asInt:    index,
		},
	})
}

func (c *compiler) VisitMul(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrMul})
}

func (c *compiler) VisitDiv(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrDiv})
}

func (c *compiler) VisitMod(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrMod})
}

func (c *compiler) VisitAdd(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrAdd})
}

func (c *compiler) VisitSub(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrSub})
}

func (c *compiler) VisitEq(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrEq})
}

func (c *compiler) VisitNe(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrNe})
}

func (c *compiler) VisitLt(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrLt})
}

func (c *compiler) VisitLe(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrLe})
}

func (c *compiler) VisitGe(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrGe})
}

func (c *compiler) VisitGt(lhs, rhs *node) {
	lhs.visit(c)
	rhs.visit(c)

	c.program.append(instr{Type: instrGt})
}

func (c *compiler) VisitHas(field, value *node) {
	field.visit(c)
	value.visit(c)

	c.program.append(instr{Type: instrHas})
}

func (c *compiler) VisitNot(n *node) {
	n.visit(c)

	c.program.append(instr{Type: instrNot})
}

func (c *compiler) VisitAnd(children ...*node) {
	for _, v := range children {
		v.visit(c)
	}

	c.program.append(instr{
		Type: instrAnd,
		Arg: argument{
			asInt: int64(len(children)),
		},
	})
}

func (c *compiler) VisitOr(children ...*node) {
	for _, v := range children {
		v.visit(c)
	}

	c.program.append(instr{
		Type: instrOr,
		Arg: argument{
			asInt: int64(len(children)),
		},
	})
}
