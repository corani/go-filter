package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	tokens  []token
	pointer int
}

func parseExpression(expression string) (*node, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		//nolint:exhaustivestruct
		return &node{
			oper:   BOOL,
			asBool: true,
		}, nil
	}

	p := &parser{tokens, 0}

	n, err := p.parseBoolExpression()
	if err != nil {
		return nil, err
	}

	if p.pointer < len(p.tokens) {
		trailing := ""

		for _, t := range p.tokens[p.pointer:] {
			if len(trailing) > 0 {
				trailing += " "
			}

			trailing += t.Value
		}

		// TODO(daniel): wrap error
		return nil, fmt.Errorf("%q", trailing)
	}

	return n, err
}

func (p *parser) token() token {
	if p.pointer < len(p.tokens) {
		return p.tokens[p.pointer]
	}

	return token{Type: INVALID, Value: ""}
}

func (p *parser) expectToken(tts ...tokenType) (token, bool) {
	if p.pointer < len(p.tokens) {
		for _, tt := range tts {
			if tt == p.tokens[p.pointer].Type {
				t := p.tokens[p.pointer]
				p.pointer++

				return t, true
			}
		}
	}

	return token{Type: INVALID, Value: ""}, false
}

func (p *parser) requireToken(tts ...tokenType) (token, error) {
	if v, ok := p.expectToken(tts...); ok {
		return v, nil
	}

	return token{}, fmt.Errorf("unexpected %v, expected on of: %v", p.token(), tts)
}

func (p *parser) parseBoolExpression() (*node, error) {
	var children []*node

	for {
		child, err := p.parseBoolTerm()
		if err != nil {
			return nil, err
		}

		children = append(children, child)

		if _, ok := p.expectToken(OR); !ok {
			break
		}
	}

	return Or(children...), nil
}

func (p *parser) parseBoolTerm() (*node, error) {
	var children []*node

	for {
		child, err := p.parseBoolFactor()
		if err != nil {
			return nil, err
		}

		children = append(children, child)

		if _, ok := p.expectToken(AND); !ok {
			break
		}
	}

	return And(children...), nil
}

func (p *parser) parseBoolFactor() (*node, error) {
	if _, ok := p.expectToken(HAS); ok {
		return p.parseHas()
	}

	if _, ok := p.expectToken(NOT); ok {
		child, err := p.parseBoolPrimary()
		if err != nil {
			return nil, err
		}

		return Not(child), nil
	}

	return p.parseBoolPrimary()
}

func (p *parser) parseHas() (*node, error) {
	lhs, err := p.parseFieldOrValue()
	if err != nil {
		return nil, err
	}

	if lhs.oper != FIELD {
		// TODO(daniel): wrap error
		return nil, fmt.Errorf("expected field, got %v", lhs.oper)
	}

	rhs, err := p.parseFieldOrValue()
	if err != nil {
		return nil, err
	}

	if rhs.oper == FIELD {
		// TODO(daniel): wrap error
		return nil, fmt.Errorf("expected value, got %v", rhs.oper)
	}

	return Has(lhs, rhs), nil
}

func (p *parser) parseBoolPrimary() (*node, error) {
	if _, ok := p.expectToken(LPAREN); !ok {
		return p.parseCondition()
	}

	result, err := p.parseBoolExpression()
	if err != nil {
		return nil, err
	}

	if _, err := p.requireToken(RPAREN); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *parser) parseCondition() (*node, error) {
	var tt tokenType

	lhs, err := p.parseArithmetic1()
	if err != nil {
		return nil, err
	}

	if t, err := p.requireToken(EQ, NE, LT, LE, GE, GT); err == nil {
		tt = t.Type
	} else {
		return nil, err
	}

	rhs, err := p.parseArithmetic1()
	if err != nil {
		return nil, err
	}

	switch tt {
	case EQ:
		return Eq(lhs, rhs), nil
	case NE:
		return Ne(lhs, rhs), nil
	case LT:
		return Lt(lhs, rhs), nil
	case LE:
		return Le(lhs, rhs), nil
	case GE:
		return Ge(lhs, rhs), nil
	case GT:
		return Gt(lhs, rhs), nil
	default:
		panic("unreachable")
	}
}

func (p *parser) parseArithmetic1() (*node, error) {
	var tt tokenType

	lhs, err := p.parseArithmetic2()
	if err != nil {
		return nil, err
	}

	if t, ok := p.expectToken(ADD, SUB); ok {
		tt = t.Type
	} else {
		return lhs, nil
	}

	rhs, err := p.parseArithmetic2()
	if err != nil {
		return nil, err
	}

	switch tt {
	case ADD:
		return Add(lhs, rhs), nil
	case SUB:
		return Sub(lhs, rhs), nil
	default:
		panic("unreachable")
	}
}

func (p *parser) parseArithmetic2() (*node, error) {
	var tt tokenType

	lhs, err := p.parseFieldOrValue()
	if err != nil {
		return nil, err
	}

	if t, ok := p.expectToken(MUL, DIV, MOD); ok {
		tt = t.Type
	} else {
		return lhs, nil
	}

	rhs, err := p.parseFieldOrValue()
	if err != nil {
		return nil, err
	}

	switch tt {
	case MUL:
		return Mul(lhs, rhs), nil
	case DIV:
		return Div(lhs, rhs), nil
	case MOD:
		return Mod(lhs, rhs), nil
	default:
		panic("unreachable")
	}
}

func (p *parser) parseFieldOrValue() (*node, error) {
	t, err := p.requireToken(FIELD, BOOL, NUMBER, STRING)
	if err != nil {
		return nil, err
	}

	switch t.Type {
	case FIELD:
		return Field(t.Value), nil
	case BOOL:
		switch strings.ToLower(t.Value) {
		case "true":
			return Bool(true), nil
		case "false":
			return Bool(false), nil
		default:
			// TODO(daniel): wrap error
			return nil, fmt.Errorf("unexpected boolean value: %v", t)
		}
	case NUMBER:
		if v, err := strconv.ParseInt(t.Value, 10, 64); err == nil {
			return Int(v), nil
		} else if v, err := strconv.ParseFloat(t.Value, 64); err == nil {
			return Float(v), nil
		} else {
			// TODO(daniel): wrap error
			return nil, fmt.Errorf("invalid number: %v", t.Value)
		}
	case STRING:
		return String(t.Value), nil
	default:
		// TODO(daniel): wrap error
		return nil, fmt.Errorf("unexpected token: %v", t)
	}
}
