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

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch {
		case r == '\\':
			if i == len(runes)-1 { // символ экранирования последний.
				return "", ErrInvalidString
			}

			if len(runes) < i+1 || (!unicode.IsDigit(runes[i+1]) && runes[i+1] != '\\') {
				return "", ErrInvalidString
			}

			if len(runes) > i+2 && unicode.IsDigit(runes[i+2]) { // если после экранирования есть кол-во повторов.
				toAdd := strings.Repeat(string(runes[i+1]), int(runes[i+2]-'0'))
				result.WriteString(toAdd)

				prev = runes[i+1]
				i = i + 2

				continue
			}

			result.WriteRune(runes[i+1])
			prev = runes[i+1]

			i++
			continue

		case unicode.IsDigit(r):
			if result.Len() == 0 { // первый символ в строке цифра.
				return "", ErrInvalidString
			}

			if i > 2 && runes[i-2] != '\\' && unicode.IsDigit(runes[i-1]) {
				return "", ErrInvalidString
			}

			if runes[i-1] != prev && i == len(runes)-1 {
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
