package compiler

import (
	"testing"

	"github.com/gscienty/causer/expr/parser"
	"github.com/gscienty/causer/runtime"
	"github.com/stretchr/testify/assert"
)

func TestCompileSimple(t *testing.T) {
	tree, err := parser.Parse("1 + 2 - 3")
	assert.Nil(t, err)
	inst, constants, err := Compile(tree)
	assert.Nil(t, err)

	expectInst := []byte{
		runtime.OpCodePush, 0x00, 0x00,
		runtime.OpCodePush, 0x00, 0x01,
		runtime.OpCodeAdd,
		runtime.OpCodePush, 0x00, 0x02,
		runtime.OpCodeSub,
	}

	assert.Equal(t, expectInst, inst)

	expectConstants := []interface{}{1, 2, 3}
	assert.Equal(t, expectConstants, constants)
}

func TestCompileCall(t *testing.T) {
	tree, err := parser.Parse("func(1 + 2 - 3) + func(1)")
	assert.Nil(t, err)
	inst, constants, err := Compile(tree)
	assert.Nil(t, err)

	expectInst := []byte{
		runtime.OpCodePush, 0x00, 0x00,
		runtime.OpCodePush, 0x00, 0x01,
		runtime.OpCodeAdd,
		runtime.OpCodePush, 0x00, 0x02,
		runtime.OpCodeSub,
		runtime.OpCodeCall, 0x00, 0x03,
		runtime.OpCodePush, 0x00, 0x00,
		runtime.OpCodeCall, 0x00, 0x03,
		runtime.OpCodeAdd,
	}

	assert.Equal(t, 4, len(constants))
	assert.Equal(t, expectInst, inst)
}
