package ast

import "reflect"

type Node interface {
	Type() reflect.Type
}

type base struct {
	nodeType reflect.Type
}

func (b *base) Type() reflect.Type { return b.nodeType }

type UnaryNode struct {
	base
	Operator string
	Expr     Node
}

type BinaryNode struct {
	base
	Operator string
	Left     Node
	Right    Node
}

type MethodNode struct {
	base
	Node      Node
	Method    string
	Arguments []Node
}

type FunctionNode struct {
	base
	Name      string
	Arguments []Node
}

type PropertyNode struct {
	base
	Node     Node
	Property string
}

type BoolNode struct {
	base
	Value bool
}

type NilNode struct {
	base
}

type IdentifierNode struct {
	base
	Value string
}

type FloatNode struct {
	base
	Value float64
}

type IntNode struct {
	base
	Value int
}

type StringNode struct {
	base
	Value string
}

type Tree struct {
	Root Node
}
