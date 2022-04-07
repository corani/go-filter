package filter

import (
	"fmt"
	"strconv"
	"strings"
)

type tokenType byte

const (
	INVALID tokenType = iota

	LPAREN
	RPAREN

	EQ
	NE
	GT
	GE
	LT
	LE

	AND
	OR
	NOT
	HAS

	FIELD
	BOOL
	NUMBER
	STRING

	ADD
	SUB
	MUL
	DIV
	MOD
)

func (tt tokenType) String() string {
	switch tt {
	case INVALID:
		return "<invalid>"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case EQ:
		return "eq"
	case NE:
		return "ne"
	case GT:
		return "gt"
	case GE:
		return "ge"
	case LT:
		return "lt"
	case LE:
		return "le"
	case AND:
		return "and"
	case OR:
		return "or"
	case NOT:
		return "not"
	case HAS:
		return "has"
	case FIELD:
		return "field"
	case BOOL:
		return "bool"
	case NUMBER:
		return "number"
	case STRING:
		return "string"
	case ADD:
		return "and"
	case SUB:
		return "sub"
	case MUL:
		return "mul"
	case DIV:
		return "div"
	case MOD:
		return "mod"
	default:
		return fmt.Sprintf("unknown: %v", byte(tt))
	}
}

type token struct {
	Type  tokenType
	Value string
}

func tokenize(expression string) ([]token, error) {
	var (
		state  int
		tokens []token
		buf    []byte
	)

	const (
		initial int = iota
		squoted
		escapeSquoted
		word
		number
		decimal
	)

	for i := 0; i < len(expression); i++ {
		ch := expression[i]

	retry:
		switch state {
		case initial:
			switch ch {
			case '(':
				tokens = append(tokens, token{LPAREN, "("})
			case ')':
				tokens = append(tokens, token{LPAREN, "("})
			case '\'':
				state = squoted
			case ' ', '\t', '\n', '\r':
				// skip whitespace
			case '.':
				buf = append(buf, ch)
				state = decimal
			default:
				switch {
				case (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z'):
					buf = append(buf, ch)
					state = word
				case (ch >= '0' && ch <= '9') || ch == '+' || ch == '-':
					buf = append(buf, ch)
					state = number
				default:
					// TODO(daniel): return error
				}
			}
		case word:
			if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
				buf = append(buf, ch)
			} else {
				if v, ok := checkIdentifier(string(buf)); ok {
					tokens = append(tokens, v)
				} else {
					tokens = append(tokens, token{FIELD, string(buf)})
				}

				buf = nil
				state = initial

				goto retry
			}
		case number:
			switch {
			case ch >= '0' && ch <= '9':
				buf = append(buf, ch)
			case ch == '.':
				buf = append(buf, ch)
				state = decimal
			default:
				if _, err := strconv.ParseInt(string(buf), 10, 64); err != nil {
					// TODO(daniel): wrap error
					return nil, err
				}

				tokens = append(tokens, token{NUMBER, string(buf)})
				buf = nil
				state = initial

				goto retry
			}
		case decimal:
			switch {
			case ch >= '0' && ch <= '9':
				buf = append(buf, ch)
			default:
				if _, err := strconv.ParseFloat(string(buf), 64); err != nil {
					// TODO(daniel): wrap error
					return nil, err
				}

				tokens = append(tokens, token{NUMBER, string(buf)})
				buf = nil
				state = initial

				goto retry
			}
		case squoted:
			switch ch {
			case '\\':
				state = escapeSquoted
			case '\'':
				tokens = append(tokens, token{STRING, string(buf)})
				buf = nil
				state = initial
			default:
				buf = append(buf, ch)
			}
		case escapeSquoted:
			buf = append(buf, ch)
			state = squoted
		}
	}

	if buf != nil {
		switch state {
		case number, decimal:
			tokens = append(tokens, token{NUMBER, string(buf)})
		case word:
			if t, ok := checkIdentifier(string(buf)); ok {
				tokens = append(tokens, t)
			} else {
				tokens = append(tokens, token{FIELD, string(buf)})
			}
		default:
			// TODO(daniel): return error
		}
	}

	return tokens, nil
}

func checkIdentifier(s string) (token, bool) {
	switch strings.ToLower(s) {
	case "eq":
		return token{EQ, s}, true
	case "ne":
		return token{NE, s}, true
	case "gt":
		return token{GT, s}, true
	case "ge":
		return token{GE, s}, true
	case "lt":
		return token{LT, s}, true
	case "le":
		return token{LE, s}, true
	case "and":
		return token{AND, s}, true
	case "or":
		return token{OR, s}, true
	case "not":
		return token{NOT, s}, true
	case "has":
		return token{HAS, s}, true
	case "add":
		return token{ADD, s}, true
	case "sub":
		return token{SUB, s}, true
	case "mul":
		return token{MUL, s}, true
	case "div":
		return token{DIV, s}, true
	case "mod":
		return token{MOD, s}, true
	case "true", "false":
		return token{BOOL, s}, true
	default:
		return token{}, false
	}
}
