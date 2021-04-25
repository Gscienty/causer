package runtime

const (
	OpCodeAdd byte = iota
	OpCodeSub
	OpCodeMul
	OpCodeDiv
	OpCodePow
	OpCodeMod
	OpCodePop
	OpCodePush
	OpCodeCall
	OpCodeNot
	OpCodeNegate
	OpCodeProperty
	OpCodeFetch
	OpCodeTrue
	OpCodeFalse
	OpCodeNil
)
