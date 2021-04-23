package parser

import (
	"testing"

	"github.com/gscienty/causer/expr/ast"
	"github.com/stretchr/testify/assert"
)

func TestParseSimple(t *testing.T) {
	root, err := Parse("a + b / 5 - (10)")
	assert.Nil(t, err)

	binaryOp, ok := root.Root.(*ast.BinaryNode)
	assert.True(t, ok)
	assert.Equal(t, "-", binaryOp.Operator)

	croot := binaryOp.Left
	binaryOp, ok = croot.(*ast.BinaryNode)
	assert.True(t, ok)
	assert.Equal(t, "+", binaryOp.Operator)

	binaryOp, ok = root.Root.(*ast.BinaryNode).Left.(*ast.BinaryNode).Right.(*ast.BinaryNode)
	assert.True(t, ok)
	assert.Equal(t, "/", binaryOp.Operator)
}
