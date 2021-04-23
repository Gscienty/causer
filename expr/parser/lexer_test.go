package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerString(t *testing.T) {
	tokens, err := Lexer(`"\"foo\""`)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, "\"foo\"", tokens[0].Value)
	assert.Equal(t, TokenKindString, tokens[0].Kind)
	assert.Equal(t, TokenKindEOF, tokens[1].Kind)
}

func TestLexerNumber(t *testing.T) {
	tokens, err := Lexer("123.45_6")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, "123.45_6", tokens[0].Value)
	assert.Equal(t, TokenKindNumber, tokens[0].Kind)
	assert.Equal(t, TokenKindEOF, tokens[1].Kind)
}
func TestLexerDotNumber(t *testing.T) {
	tokens, err := Lexer(".12345_6")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tokens))
	assert.Equal(t, ".12345_6", tokens[0].Value)
	assert.Equal(t, TokenKindNumber, tokens[0].Kind)
	assert.Equal(t, TokenKindEOF, tokens[1].Kind)
}

func TestLexerBracket(t *testing.T) {
	tokens, err := Lexer("(123.45_6)")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(tokens))
	assert.Equal(t, TokenKindBracket, tokens[0].Kind)
	assert.Equal(t, TokenKindNumber, tokens[1].Kind)
	assert.Equal(t, TokenKindBracket, tokens[2].Kind)
	assert.Equal(t, TokenKindEOF, tokens[3].Kind)
	assert.Equal(t, "(", tokens[0].Value)
	assert.Equal(t, "123.45_6", tokens[1].Value)
	assert.Equal(t, ")", tokens[2].Value)
}

func TestLexerKeywords(t *testing.T) {
	tokens, err := Lexer("var in [1, 2, 3, 4]")
	assert.Nil(t, err)
	assert.Equal(t, 12, len(tokens))
	assert.Equal(t, TokenKindIdentifier, tokens[0].Kind)
	assert.Equal(t, TokenKindOperator, tokens[1].Kind)
	assert.Equal(t, TokenKindBracket, tokens[2].Kind)
	assert.Equal(t, TokenKindNumber, tokens[3].Kind)
	assert.Equal(t, TokenKindOperator, tokens[4].Kind)
	assert.Equal(t, TokenKindNumber, tokens[5].Kind)
	assert.Equal(t, TokenKindOperator, tokens[6].Kind)
	assert.Equal(t, TokenKindNumber, tokens[7].Kind)
	assert.Equal(t, TokenKindOperator, tokens[8].Kind)
	assert.Equal(t, TokenKindNumber, tokens[9].Kind)
	assert.Equal(t, TokenKindBracket, tokens[10].Kind)
	assert.Equal(t, TokenKindEOF, tokens[11].Kind)
}

func TestLexerOperator(t *testing.T) {
	token, err := Lexer("1+2/3*4")
	assert.Nil(t, err)
	assert.Equal(t, 8, len(token))
	assert.Equal(t, TokenKindNumber, token[0].Kind)
	assert.Equal(t, TokenKindOperator, token[1].Kind)
	assert.Equal(t, TokenKindNumber, token[2].Kind)
	assert.Equal(t, TokenKindOperator, token[3].Kind)
	assert.Equal(t, TokenKindNumber, token[4].Kind)
	assert.Equal(t, TokenKindOperator, token[5].Kind)
	assert.Equal(t, TokenKindNumber, token[6].Kind)
	assert.Equal(t, TokenKindEOF, token[7].Kind)
}
