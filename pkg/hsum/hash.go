package hsum

import (
	"fmt"
)

type HashInterface interface {
	HashByString(string) string
	Hash([]byte) string
}

func New() HashInterface {
	return &fnvImpl{}
}

type fnvImpl struct{}

func (f *fnvImpl) HashByString(str string) string {
	return f.Hash([]byte(str))
}

// Hash https://ru.wikipedia.org/wiki/FNV
func (f *fnvImpl) Hash(input []byte) string {
	var result uint64 = 0xcbf29ce484222325
	for _, char := range input {
		result ^= uint64(char)
		result *= 0x00000100000001b3
	}

	return fmt.Sprintf("%016x", result)
}
