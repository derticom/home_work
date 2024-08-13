package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	runes := []rune(input)

	var prev string

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '\\':
			if len(runes) < i+1 || (!unicode.IsDigit(runes[i+1]) && runes[i+1] != '\\') {
				return "", ErrInvalidString
			}

			result.WriteRune(runes[i+1])
			prev = string(runes[i+1])

			i++
			continue
		case unicode.IsDigit(r):
			if result.Len() == 0 && len(prev) == 0 {
				return "", ErrInvalidString
			}

			if i > 2 && runes[i-2] != '\\' && unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}

			for i := 1; i < int(r-'0'); i++ {
				result.WriteString(prev)
			}
		case unicode.IsLetter(r):
			if (len(runes) > i+1 && runes[i+1]-'0' != 0) || i == len(runes)-1 {
				result.WriteRune(r)
			}
		default:
			return "", ErrInvalidString
		}

		prev = string(r)

		fmt.Println(i, prev)
	}

	return result.String(), nil
}
