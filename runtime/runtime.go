package runtime

import (
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

type Runtime struct {
	stack        []interface{}
	constants    []interface{}
	instructions []byte

	instructionPointer int

	instFunc map[byte]func() error
	instImpl map[string][]interface{}

	env interface{}
}

const (
	runtimeOpAdd = "+"
	runtimeOpSub = "-"
	runtimeOpMul = "*"
	runtimeOpDiv = "/"
	runtimeOpMod = "%"
	runtimeOpPow = "^"
)

func New(instructions []byte, constants []interface{}, env interface{}) *Runtime {
	rt := &Runtime{
		stack:        make([]interface{}, 0),
		constants:    constants,
		instructions: instructions,
		instImpl:     make(map[string][]interface{}),
		env:          env,
	}

	rt.instFunc = map[byte]func() error{
		OpCodeAdd:      rt.instBinaryOp(runtimeOpAdd),
		OpCodeSub:      rt.instBinaryOp(runtimeOpSub),
		OpCodeMul:      rt.instBinaryOp(runtimeOpMul),
		OpCodeDiv:      rt.instBinaryOp(runtimeOpDiv),
		OpCodeMod:      rt.instBinaryOp(runtimeOpMod),
		OpCodePow:      rt.instBinaryOp(runtimeOpPow),
		OpCodePop:      rt.instPop,
		OpCodePush:     rt.instPush,
		OpCodeCall:     rt.instCall,
		OpCodeFetch:    rt.instFetch,
		OpCodeProperty: rt.instProperty,
	}

	return rt
}

func (r *Runtime) readArg() uint16 {
	ret := binary.BigEndian.Uint16(r.instructions[r.instructionPointer : r.instructionPointer+2])
	r.instructionPointer += 2
	return ret
}

func (r *Runtime) readConstant() interface{} { return r.constants[r.readArg()] }

func (r *Runtime) push(v interface{}) { r.stack = append(r.stack, v) }
func (r *Runtime) pop() interface{} {
	v := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return v
}

func (r *Runtime) Run() (interface{}, error) {
	for r.instructionPointer < len(r.instructions) {
		op := r.instructions[r.instructionPointer]
		r.instructionPointer++

		if instFunc, ok := r.instFunc[op]; ok {
			instFunc()
		} else {
			return nil, fmt.Errorf("unexcepted instruction")
		}
	}

	return r.pop(), nil
}

func (r *Runtime) Register(name string, fn interface{}) {
	if _, ok := r.instImpl[name]; !ok {
		r.instImpl[name] = make([]interface{}, 0)
	}

	r.instImpl[name] = append(r.instImpl[name], fn)
}

func (r *Runtime) instBinaryOp(op string) func() error {
	return func() error {
		right := r.pop()
		left := r.pop()

		impls, ok := r.instImpl[op]
		if !ok {
			return fmt.Errorf("invalid operator %s", op)
		}
		for _, impl := range impls {
			if reflect.TypeOf(impl).Kind() == reflect.Func {
				implType := reflect.TypeOf(impl)
				if implType.NumIn() != 2 && implType.NumOut() != 1 {
					continue
				}
				if implType.In(0).Kind() != reflect.TypeOf(left).Kind() || implType.In(1).Kind() != reflect.TypeOf(right).Kind() {
					continue
				}

				implValue := reflect.ValueOf(impl)
				ret := implValue.Call([]reflect.Value{reflect.ValueOf(left), reflect.ValueOf(right)})

				r.push(ret[0].Interface())
				return nil
			}
		}
		return fmt.Errorf("invalid operator %s", op)
	}
}

func (r *Runtime) instPop() error {
	r.pop()
	return nil
}

func (r *Runtime) instPush() error {
	r.push(r.readConstant())
	return nil
}

func (r *Runtime) instProperty() error {
	instance := r.pop()
	prop := r.readConstant()
	r.push(r.fetch(instance, prop))
	return nil
}

func (r *Runtime) instFetch() error {
	r.push(r.fetch(r.env, r.readConstant()))
	return nil
}

func (r *Runtime) instNil() error {

	return nil
}

func (r *Runtime) instCall() error {
	call := r.readConstant().(Call)
	in := make([]reflect.Value, call.ArgumentsCnt)
	for i := call.ArgumentsCnt; i > 0; i-- {
		param := r.pop()
		if param == nil && reflect.TypeOf(param) == nil {
			in[i-1] = reflect.ValueOf(&param).Elem()
		} else {
			in[i-1] = reflect.ValueOf(param)
		}
	}
	r.push((*r.fetchFn(call.Name)).Call(in)[0].Interface())
	return nil
}

func (r *Runtime) fetchFn(name string) *reflect.Value {
	v := reflect.ValueOf(r.env)

	if v.NumMethod() > 0 {
		method := v.MethodByName(name)
		if method.IsValid() {
			return &method
		}
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		ret := v.MapIndex(reflect.ValueOf(name))
		if ret.IsValid() && ret.CanInterface() {
			ret = ret.Elem()
			return &ret
		}
	case reflect.Struct:
		ret := v.FieldByName(name)
		if ret.IsValid() {
			return &ret
		}
	}

	return nil
}

func (r *Runtime) fetch(env interface{}, identifiy interface{}) interface{} {
	envValue := reflect.ValueOf(env)

	if envValue.Kind() == reflect.Ptr && reflect.Indirect(envValue).Kind() == reflect.Struct {
		envValue = reflect.Indirect(envValue)
	}

	switch envValue.Kind() {
	case reflect.Array, reflect.Slice, reflect.String:
		v := envValue.Index(cast.ToInt(identifiy))
		if v.IsValid() && v.CanInterface() {
			return v.Interface()
		}

	case reflect.Map:
		v := envValue.MapIndex(reflect.ValueOf(identifiy))
		if v.IsValid() {
			if v.CanInterface() {
				return v.Interface()
			} else {
				return reflect.Zero(reflect.TypeOf(env).Elem()).Interface()
			}
		}

	case reflect.Struct:
		v := envValue.FieldByName(reflect.ValueOf(identifiy).String())
		if v.IsValid() && v.CanInterface() {
			return v.Interface()
		}
	}

	return nil
}
