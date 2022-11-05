package bytesize

import (
	"strconv"
	"strings"
	"unicode"
)

const (
	B uint64 = 1
	K uint64 = 1 << (10 * iota)
	M
	G
)

var unitMap = map[string]uint64{
	"B": B,
	"K": K,
	"M": M,
	"G": G,
}

func Parse(input string) uint64 {
	input = strings.TrimSpace(input)
	split := make([]string, 0)
	for i, char := range input {
		if !unicode.IsDigit(char) {
			split = append(split, strings.TrimSpace(input[:i]))
			split = append(split, strings.TrimSpace(input[i:]))
			break
		}
	}

	if len(split) != 2 {
		return 0
	}

	unit, ok := unitMap[strings.ToUpper(split[1])]
	if !ok {
		return 0
	}

	value, err := strconv.ParseUint(split[0], 10, 64)
	if err != nil {
		return 0
	}

	return value * unit
}
