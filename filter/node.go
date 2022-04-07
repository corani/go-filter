package filter

type node struct {
	oper     tokenType
	children []*node
	asBool   bool
	asInt    int64
	asFloat  float64
	asString string
	asField  string
}

func (n *node) visit(v Visitor) {
	switch n.oper {
	case BOOL:
		v.VisitBool(n)
	case NUMBER:
		v.VisitNumber(n)
	case STRING:
		v.VisitString(n)
	case FIELD:
		v.VisitField(n)
	case MUL:
		v.VisitMul(n.children[0], n.children[1])
	case DIV:
		v.VisitDiv(n.children[0], n.children[1])
	case MOD:
		v.VisitMod(n.children[0], n.children[1])
	case ADD:
		v.VisitAdd(n.children[0], n.children[1])
	case SUB:
		v.VisitSub(n.children[0], n.children[1])
	case EQ:
		v.VisitEq(n.children[0], n.children[1])
	case NE:
		v.VisitNe(n.children[0], n.children[1])
	case LT:
		v.VisitLt(n.children[0], n.children[1])
	case LE:
		v.VisitLe(n.children[0], n.children[1])
	case GE:
		v.VisitGe(n.children[0], n.children[1])
	case GT:
		v.VisitGt(n.children[0], n.children[1])
	case HAS:
		v.VisitHas(n.children[0], n.children[1])
	case NOT:
		v.VisitNot(n.children[0])
	case AND:
		v.VisitAnd(n.children...)
	case OR:
		v.VisitOr(n.children...)
	}
}

func Bool(v bool) *node {
	return &node{
		oper:   BOOL,
		asBool: v,
	}
}

func Int(v int64) *node {
	return &node{
		oper:    NUMBER,
		asInt:   v,
		asFloat: float64(v),
	}
}

func Float(v float64) *node {
	return &node{
		oper:    NUMBER,
		asInt:   int64(v),
		asFloat: v,
	}
}

func String(v string) *node {
	return &node{
		oper:     STRING,
		asString: v,
	}
}

func Field(v string) *node {
	return &node{
		oper:    FIELD,
		asField: v,
	}
}

func Mul(lhs, rhs *node) *node {
	return &node{
		oper:     MUL,
		children: []*node{lhs, rhs},
	}
}

func Div(lhs, rhs *node) *node {
	return &node{
		oper:     DIV,
		children: []*node{lhs, rhs},
	}
}

func Mod(lhs, rhs *node) *node {
	return &node{
		oper:     MOD,
		children: []*node{lhs, rhs},
	}
}

func Add(lhs, rhs *node) *node {
	return &node{
		oper:     ADD,
		children: []*node{lhs, rhs},
	}
}

func Sub(lhs, rhs *node) *node {
	return &node{
		oper:     SUB,
		children: []*node{lhs, rhs},
	}
}

func Eq(lhs, rhs *node) *node {
	return &node{
		oper:     EQ,
		children: []*node{lhs, rhs},
	}
}

func Ne(lhs, rhs *node) *node {
	return &node{
		oper:     NE,
		children: []*node{lhs, rhs},
	}
}

func Lt(lhs, rhs *node) *node {
	return &node{
		oper:     LT,
		children: []*node{lhs, rhs},
	}
}

func Le(lhs, rhs *node) *node {
	return &node{
		oper:     LE,
		children: []*node{lhs, rhs},
	}
}

func Ge(lhs, rhs *node) *node {
	return &node{
		oper:     GE,
		children: []*node{lhs, rhs},
	}
}

func Gt(lhs, rhs *node) *node {
	return &node{
		oper:     GT,
		children: []*node{lhs, rhs},
	}
}

func Has(field, value *node) *node {
	return &node{
		oper:     HAS,
		children: []*node{field, value},
	}
}

func Not(n *node) *node {
	return &node{
		oper:     NOT,
		children: []*node{n},
	}
}

func And(children ...*node) *node {
	if len(children) == 1 {
		return children[0]
	}

	return &node{
		oper:     AND,
		children: children,
	}
}

func Or(children ...*node) *node {
	if len(children) == 1 {
		return children[0]
	}

	return &node{
		oper:     OR,
		children: children,
	}
}
