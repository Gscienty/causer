package compiler

import (
	"reflect"

	"github.com/gscienty/causer/expr/ast"
	"github.com/gscienty/causer/runtime"
)

func Compile(tree *ast.Tree) ([]byte, []interface{}, error) {
	c := compiler{
		instructions:   make([]byte, 0),
		constants:      make([]interface{}, 0),
		constantsIndex: make(map[interface{}]uint16),
	}

	c.compile(tree.Root)

	return c.instructions, c.constants, nil
}

type compiler struct {
	instructions []byte

	constants      []interface{}
	constantsIndex map[interface{}]uint16
}

func (c *compiler) compile(node ast.Node) {
	switch n := node.(type) {
	case *ast.UnaryNode:
		c.compileUnaryNode(n)
	case *ast.BinaryNode:
		c.compileBinaryNode(n)
	case *ast.MethodNode:
		c.compileMethodNode(n)
	case *ast.FunctionNode:
		c.compileFunctionNode(n)
	case *ast.PropertyNode:
		c.compilePropertyNode(n)
	case *ast.IdentifierNode:
		c.compileIdentifierNode(n)
	case *ast.BoolNode:
		c.compileBoolNode(n)
	case *ast.NilNode:
		c.compileNilNode(n)
	case *ast.FloatNode:
		c.compileFloatNode(n)
	case *ast.IntNode:
		c.compileIntNode(n)
	case *ast.StringNode:
		c.compileStringNode(n)
	}
}

func (c *compiler) appendInstruction(instruction byte, operands ...byte) {
	c.instructions = append(c.instructions, instruction)
	c.instructions = append(c.instructions, operands...)
}

func (c *compiler) compileUnaryNode(n *ast.UnaryNode) {
	c.compile(n.Expr)

	switch n.Operator {
	case "!", "not":
		c.appendInstruction(runtime.OpCodeNot)
	case "-":
		c.appendInstruction(runtime.OpCodeNegate)
	}
}

func (c *compiler) compileBinaryNode(n *ast.BinaryNode) {
	c.compile(n.Left)
	c.compile(n.Right)

	switch n.Operator {
	case "+":
		c.appendInstruction(runtime.OpCodeAdd)
	case "-":
		c.appendInstruction(runtime.OpCodeSub)
	case "*":
		c.appendInstruction(runtime.OpCodeMul)
	case "/":
		c.appendInstruction(runtime.OpCodeDiv)
	case "^":
		c.appendInstruction(runtime.OpCodePow)
	case "%":
		c.appendInstruction(runtime.OpCodeMod)
	}
}

func (c *compiler) compileMethodNode(n *ast.MethodNode) {
	c.compile(n.Node)
	for _, arg := range n.Arguments {
		c.compile(arg)
	}

	c.appendInstruction(runtime.OpCodeCall, c.newConstant(runtime.Call{Name: n.Method, ArgumentsCnt: len(n.Arguments)})...)
}

func (c *compiler) compileFunctionNode(n *ast.FunctionNode) {
	for _, arg := range n.Arguments {
		c.compile(arg)
	}

	c.appendInstruction(runtime.OpCodeCall, c.newConstant(runtime.Call{Name: n.Name, ArgumentsCnt: len(n.Arguments)})...)
}

func (c *compiler) compilePropertyNode(n *ast.PropertyNode) {
	c.compile(n.Node)
	c.appendInstruction(runtime.OpCodeProperty, c.newConstant(n.Property)...)
}

func (c *compiler) compileIdentifierNode(n *ast.IdentifierNode) {
	c.appendInstruction(runtime.OpCodeFetch, c.newConstant(n.Value)...)
}

func (c *compiler) compileBoolNode(n *ast.BoolNode) {
	if n.Value {
		c.appendInstruction(runtime.OpCodeTrue)
	} else {
		c.appendInstruction(runtime.OpCodeFalse)
	}
}

func (c *compiler) compileFloatNode(n *ast.FloatNode) {
	c.appendInstruction(runtime.OpCodePush, c.newConstant(n.Value)...)
}

func (c *compiler) compileIntNode(n *ast.IntNode) {
	c.appendInstruction(runtime.OpCodePush, c.newConstant(n.Value)...)
}

func (c *compiler) compileNilNode(n *ast.NilNode) {
	c.appendInstruction(runtime.OpCodeNil)
}

func (c *compiler) compileStringNode(n *ast.StringNode) {
	c.appendInstruction(runtime.OpCodePush, c.newConstant(n.Value)...)
}

func (c *compiler) newConstant(i interface{}) []byte {
	hashable := true
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice, reflect.Map:
		hashable = false
	}

	if hashable {
		if ret, ok := c.constantsIndex[i]; ok {
			return encode(ret)
		}
	}

	c.constants = append(c.constants, i)
	if hashable {
		c.constantsIndex[i] = uint16(len(c.constants) - 1)
	}

	return encode(uint16(len(c.constants) - 1))
}
