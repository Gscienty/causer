package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type lexer struct {
	source   string
	tokens   []Token
	start    int
	end      int
	width    int
	startPos Position
	locPos   Position
	prevPos  Position
	err      error
}

type lexerStateFunc func(l *lexer) lexerStateFunc

const eof rune = -1

var keywords []string = []string{
	"not",
	"in",
	"and",
	"or",
}

func (l *lexer) nextAlpha() rune {
	alpha, width := l.peekAlpha()
	if alpha == eof {
		return eof
	}

	l.width = width
	l.end += width

	l.prevPos = l.locPos
	if alpha == '\n' {
		l.locPos.Line++
		l.locPos.Offset = 0
	} else {
		l.locPos.Offset++
	}

	return alpha
}

func (l *lexer) prevAlpha() {
	l.end -= l.width
	l.locPos = l.prevPos
}

func (l *lexer) peekAlpha() (rune, int) {
	if l.end >= len(l.source) {
		l.width = 0
		return eof, 0
	}

	return utf8.DecodeRuneInString(l.source[l.end:])
}

func (l *lexer) ignore() {
	l.start = l.end
	l.startPos = l.locPos
}

func (l *lexer) word() string {
	return l.source[l.start:l.end]
}

func (l *lexer) scanEscape(quote rune) rune {
	alpha := l.nextAlpha()
	switch alpha {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		alpha = l.nextAlpha()
	default:
		l.err = fmt.Errorf("unexpected escape")
	}

	return alpha
}

func (l *lexer) scanString(quote rune) {
	alpha := l.nextAlpha()
	for alpha != quote {
		if alpha == '\n' || alpha == eof {
			l.err = fmt.Errorf("unexpected terminated")
			return
		}

		if alpha == '\\' {
			alpha = l.scanEscape(quote)
		} else {
			alpha = l.nextAlpha()
		}
	}
}

func (l *lexer) accept(valid string) bool {
	alpha, _ := l.peekAlpha()
	if strings.ContainsRune(valid, alpha) {
		l.nextAlpha()
		return true
	}
	return false
}

func (l *lexer) scanNumber() bool {
	dig := "0123456789_"
	if l.accept("0") {
		if l.accept("xX") {
			dig = "0123456789abcdefABCDEF_"
		} else if l.accept("oO") {
			dig = "01234567_"
		} else if l.accept("bB") {
			dig = "01_"
		}
	}

	next := l.accept(dig)
	for next {
		next = l.accept(dig)
	}

	if l.accept(".") {
		next := l.accept(dig)
		for next {
			next = l.accept(dig)
		}
	}

	return true
}

func (l *lexer) productEOF() {
	l.tokens = append(l.tokens, Token{
		Kind:     TokenKindEOF,
		Position: l.prevPos,
	})

	l.start = l.end
}

func (l *lexer) product(kind Kind, word string) {
	l.tokens = append(l.tokens, Token{
		Kind:     kind,
		Position: l.startPos,
		Value:    word,
	})

	l.start = l.end
	l.startPos = l.locPos
}

func rootState(l *lexer) lexerStateFunc {
	switch alpha := l.nextAlpha(); {
	case alpha == eof:
		l.productEOF()
		return nil
	case isSpace(alpha):
		l.ignore()
	case alpha == '\'' || alpha == '"':
		l.scanString(alpha)
		word, err := unescape(l.word())
		if err != nil {
			l.err = err
		}
		l.product(TokenKindString, word)
	case '0' <= alpha && alpha <= '9':
		l.prevAlpha()
		return numberState
	case strings.ContainsRune("([{", alpha):
		l.product(TokenKindBracket, l.word())
	case strings.ContainsRune(")]}", alpha):
		l.product(TokenKindBracket, l.word())
	case strings.ContainsRune("?:%,+-*/^", alpha):
		l.product(TokenKindOperator, l.word())
	case strings.ContainsRune("&|!=<>", alpha):
		l.accept("&|=")
		l.product(TokenKindOperator, l.word())
	case alpha == '.':
		l.prevAlpha()
		return dotState
	case isAlphaNumeric(alpha):
		l.prevAlpha()
		return identifierState
	}

	return rootState
}

func identifierState(l *lexer) lexerStateFunc {
loop:
	for {
		switch alpha := l.nextAlpha(); {
		case isAlphaNumeric(alpha):
		default:
			l.prevAlpha()

			word := l.word()
			isOperator := false
			for _, op := range keywords {
				if word == op {
					l.product(TokenKindOperator, word)
					isOperator = true
					break
				}
			}

			if !isOperator {
				l.product(TokenKindIdentifier, word)
			}

			break loop
		}
	}

	return rootState
}

func dotState(l *lexer) lexerStateFunc {
	l.nextAlpha()
	alpha, _ := l.peekAlpha()
	if strings.ContainsRune("0123456789", alpha) {
		return numberState
	}
	l.product(TokenKindOperator, l.word())
	return rootState
}

func numberState(l *lexer) lexerStateFunc {
	if !l.scanNumber() {
		l.err = fmt.Errorf("bad number syntax: %q", l.word())
		return nil
	}

	l.product(TokenKindNumber, l.word())
	return rootState
}

func Lexer(source string) ([]Token, error) {
	l := &lexer{
		source:   source,
		tokens:   make([]Token, 0),
		startPos: Position{Line: 1, Offset: 0},
		locPos:   Position{Line: 1, Offset: 0},
		prevPos:  Position{Line: 1, Offset: 0},
	}

	for state := rootState; state != nil && l.err == nil; {
		state = state(l)
	}

	if l.err != nil {
		return nil, l.err
	}
	return l.tokens, nil
}
