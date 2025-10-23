package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var result strings.Builder
	runes := []rune(str)
	length := len(runes)

	if length == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < length; i++ {
		current := runes[i]

		if current == '\\' {
			if i+1 >= length {
				return "", ErrInvalidString
			}
			i++
			current = runes[i]

			if i+1 < length && unicode.IsDigit(runes[i+1]) {
				count, _ := strconv.Atoi(string(runes[i+1]))
				result.WriteString(strings.Repeat(string(current), count))
				i++
			} else {
				result.WriteRune(current)
			}
			continue
		}

		if unicode.IsDigit(current) {
			return "", ErrInvalidString
		}

		if i+1 < length && unicode.IsDigit(runes[i+1]) {
			count, _ := strconv.Atoi(string(runes[i+1]))
			result.WriteString(strings.Repeat(string(current), count))
			i++
		} else {
			result.WriteRune(current)
		}
	}

	return result.String(), nil
}
