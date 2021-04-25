package runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuntimeAdd(t *testing.T) {
	r := New([]byte{
		OpCodePush, 0x00, 0x00,
		OpCodePush, 0x00, 0x01,
		OpCodeAdd,
		OpCodePush, 0x00, 0x02,
		OpCodeSub,
		OpCodePush, 0x00, 0x02,
		OpCodeCall, 0x00, 0x03,
	}, []interface{}{8, 7.3, 10, Call{Name: "mul", ArgumentsCnt: 2}},
		map[string]interface{}{
			"mul": func(a float64, b int) float64 { return a * float64(b) },
		},
	)

	r.Register("+", func(left int, right float64) float64 { return float64(left) + right })
	r.Register("+", func(left int, right int) int { return left + right })
	r.Register("-", func(left float64, right int) float64 { return left - float64(right) })

	ret, err := r.Run()
	assert.Nil(t, err)

	fmt.Printf("%v", ret)
}
