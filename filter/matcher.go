package filter

import "fmt"

type valueKind byte

const (
	valueKindBool valueKind = iota
	valueKindNumber
	valueKindString
)

func (k valueKind) String() string {
	switch k {
	case valueKindBool:
		return "bool"
	case valueKindNumber:
		return "number"
	case valueKindString:
		return "string"
	}

	return fmt.Sprintf("unknown(%d)", byte(k))
}

type value struct {
	kind     valueKind
	asBool   bool
	asFloat  float64
	asString string
}

type stack struct {
	values  []value
	pointer int
}

func (s *stack) push(v value) {
	if s.pointer < len(s.values)-1 {
		s.values[s.pointer] = v
	} else {
		s.values = append(s.values, v)
	}

	s.pointer++
}

func (s *stack) pop() value {
	s.pointer--
	return s.values[s.pointer]
}
