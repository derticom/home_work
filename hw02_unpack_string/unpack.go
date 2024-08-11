package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder

	runes := []rune(input)

	var prev rune

	for i, r := range runes {
		switch {
		case unicode.IsDigit(r):
			if (result.Len() == 0 && prev == 0) || unicode.IsDigit(prev) {
				return "", ErrInvalidString
			}

			for i := 1; i < int(r-'0'); i++ {
				result.WriteRune(prev)
			}
		case unicode.IsLetter(r):
			if (len(runes) > i+1 && runes[i+1]-'0' != 0) || i == len(runes)-1 {
				result.WriteRune(r)
			}
		default:
			return "", ErrInvalidString
		}

		prev = r
	}

	return result.String(), nil
}
