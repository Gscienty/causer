package parser

type Kind string

const (
	TokenKindOperator   Kind = "operand"
	TokenKindNumber     Kind = "number"
	TokenKindString     Kind = "string"
	TokenKindBracket    Kind = "bracket"
	TokenKindIdentifier Kind = "identifier"
	TokenKindEOF        Kind = "eof"
)

type Token struct {
	Position Position
	Kind     Kind
	Value    string
}
