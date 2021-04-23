package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func isSpace(alpha rune) bool { return unicode.IsSpace(alpha) }

func isAlphaNumeric(alpha rune) bool {
	return alpha == '_' || alpha == '$' || unicode.IsLetter(alpha) || unicode.IsDigit(alpha)
}

var newline = strings.NewReplacer("\r\n", "\n", "\r", "\n")

func unescape(input string) (string, error) {
	input = newline.Replace(input)

	if len(input) < 2 {
		return input, fmt.Errorf("unable escape string")
	}
	if input[0] != input[len(input)-1] || (input[0] != '\'' && input[0] != '"') {
		return input, fmt.Errorf("unable escape string")
	}

	var tmp [utf8.UTFMax]byte
	buf := make([]byte, 0, 3*len(input)/2)

	input = input[1 : len(input)-1]
	for len(input) > 0 {
		alpha, multibyte, rest, err := unescapeAlpha(input)
		if err != nil {
			return "", err
		}
		input = rest
		if alpha < utf8.RuneSelf || !multibyte {
			buf = append(buf, byte(alpha))
		} else {
			n := utf8.EncodeRune(tmp[:], alpha)
			buf = append(buf, tmp[:n]...)
		}
	}

	return string(buf), nil
}

func unescapeAlpha(input string) (rune, bool, string, error) {
	switch c := input[0]; {
	case c > utf8.RuneSelf:
		alpha, size := utf8.DecodeLastRuneInString(input)
		return alpha, true, input[size:], nil
	case c != '\\':
		return rune(c), false, input[1:], nil
	}

	if len(input) <= 1 {
		return 0, false, "", fmt.Errorf("unable escape string '\\' as last character")
	}

	alpha := input[1]
	input = input[2:]

	var value rune
	switch alpha {
	case 'a':
		value = '\a'
	case 'b':
		value = '\b'
	case 'f':
		value = '\f'
	case 'n':
		value = '\n'
	case 'r':
		value = '\r'
	case 't':
		value = '\t'
	case 'v':
		value = '\v'
	case '\\':
		value = '\\'
	case '\'':
		value = '\''
	case '"':
		value = '"'
	default:
		return 0, false, "", fmt.Errorf("unable escape \\%c", alpha)
	}

	return value, false, input, nil
}
