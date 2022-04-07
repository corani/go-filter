package filter

import (
	"fmt"
	"reflect"
)

type instructionType byte

const (
	instPushBool instructionType = iota
	instrPushNumber
	instrPushString
	instrPushField
	instrMul
	instrDiv
	instrMod
	instrAdd
	instrSub
	instrEq
	instrNe
	instrLt
	instrLe
	instrGe
	instrGt
	instrHas
	instrNot
	instrAnd
	instrOr
)

type argument struct {
	asBool   bool
	asInt    int64
	asFloat  float64
	asString string
}

type instr struct {
	Type instructionType
	Arg  argument
}

type program struct {
	instructions []instr
	element      reflect.Type
}

func (p *program) append(i instr) {
	p.instructions = append(p.instructions, i)
}

func (p program) run(data reflect.Value) (bool, error) {
	if data.Kind() != reflect.Struct {
		return false, fmt.Errorf("unexpected type: %v", reflect.TypeOf(data).Kind())
	}

	if !data.Type().AssignableTo(p.element) {
		return false, fmt.Errorf("element type %v is not assignable to %v", data.Kind(), p.element.Kind())
	}

	s := new(stack)

	for _, instr := range p.instructions {
		switch instr.Type {
		case instPushBool:
			s.push(value{
				kind:   valueKindBool,
				asBool: instr.Arg.asBool,
			})
		case instrPushNumber:
			s.push(value{
				kind:    valueKindNumber,
				asFloat: instr.Arg.asFloat,
			})
		case instrPushString:
			s.push(value{
				kind:     valueKindString,
				asString: instr.Arg.asString,
			})
		case instrPushField:
			field := data.Field(int(instr.Arg.asInt))

			switch field.Kind() {
			case reflect.String:
				s.push(value{
					kind:     valueKindString,
					asString: field.String(),
				})
			case reflect.Bool:
				s.push(value{
					kind:   valueKindBool,
					asBool: field.Bool(),
				})
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				s.push(value{
					kind:    valueKindNumber,
					asFloat: float64(field.Int()),
				})
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				// TODO(daniel): check for overflow.
				s.push(value{
					kind:    valueKindNumber,
					asFloat: float64(field.Uint()),
				})
			case reflect.Float32, reflect.Float64:
				s.push(value{
					kind:    valueKindNumber,
					asFloat: field.Float(),
				})
			default:
				return false, fmt.Errorf("unsupported field type: %v", field.Kind())
			}
		case instrMul:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != valueKindNumber || rhs.kind != valueKindNumber {
				return false, fmt.Errorf("unsupported type in `mul`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			s.push(value{
				kind:    valueKindNumber,
				asFloat: lhs.asFloat * rhs.asFloat,
			})
		case instrDiv:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != valueKindNumber || rhs.kind != valueKindNumber {
				return false, fmt.Errorf("unsupported type in `div`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			if rhs.asFloat == 0 {
				s.push(value{
					kind:    valueKindNumber,
					asFloat: 0,
				})
			} else {
				s.push(value{
					kind:    valueKindNumber,
					asFloat: lhs.asFloat / rhs.asFloat,
				})
			}
		case instrMod:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != valueKindNumber || rhs.kind != valueKindNumber {
				return false, fmt.Errorf("unsupported type in `mod`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			s.push(value{
				kind:    valueKindNumber,
				asFloat: float64(int64(lhs.asFloat) % int64(rhs.asFloat)),
			})
		case instrAdd:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != valueKindNumber || rhs.kind != valueKindNumber {
				return false, fmt.Errorf("unsupported type in `add`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			s.push(value{
				kind:    valueKindNumber,
				asFloat: lhs.asFloat + rhs.asFloat,
			})
		case instrSub:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != valueKindNumber || rhs.kind != valueKindNumber {
				return false, fmt.Errorf("unsupported type in `sub`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			// TODO(daniel): check the order of arguments.
			s.push(value{
				kind:    valueKindNumber,
				asFloat: lhs.asFloat - rhs.asFloat,
			})
		case instrEq:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `eq`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindBool:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asBool == rhs.asBool,
				})
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat == rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString == rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrNe:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `ne`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindBool:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asBool != rhs.asBool,
				})
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat != rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString != rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrLt:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `lt`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat < rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString < rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrLe:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `le`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat <= rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString <= rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrGe:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `ge`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat >= rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString >= rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrGt:
			rhs := s.pop()
			lhs := s.pop()

			if lhs.kind != rhs.kind {
				return false, fmt.Errorf("unsupported type in `gt`: lhs=%s, rhs=%s", lhs.kind, rhs.kind)
			}

			switch lhs.kind {
			case valueKindNumber:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asFloat > rhs.asFloat,
				})
			case valueKindString:
				s.push(value{
					kind:   valueKindBool,
					asBool: lhs.asString > rhs.asString,
				})
			default:
				// TODO(daniel): report error.
			}
		case instrHas:
		// TODO(daniel): implementation.
		case instrNot:
			value := s.pop()

			if value.kind != valueKindBool {
				return false, fmt.Errorf("unsupported type in `not`: %s", value.kind)
			}

			value.asBool = !value.asBool

			s.push(value)
		case instrAnd:
			result := true

			for i := 0; i < int(instr.Arg.asInt); i++ {
				value := s.pop()

				if value.kind != valueKindBool {
					return false, fmt.Errorf("unsupported type in `or`: %s", value.kind)
				}

				result = result && value.asBool
			}

			s.push(value{
				kind:   valueKindBool,
				asBool: result,
			})
		case instrOr:
			result := false

			for i := 0; i < int(instr.Arg.asInt); i++ {
				value := s.pop()

				if value.kind != valueKindBool {
					return false, fmt.Errorf("unsupported type in `and`: %s", value.kind)
				}

				result = result || value.asBool
			}

			s.push(value{
				kind:   valueKindBool,
				asBool: result,
			})
		}
	}

	top := s.pop()

	if s.pointer > 0 || top.kind != valueKindBool {
		return false, fmt.Errorf("unexpected stack after run: len=%d, top=%v", s.pointer, top)
	}

	return top.asBool, nil
}
