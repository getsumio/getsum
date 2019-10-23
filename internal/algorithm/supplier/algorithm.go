package supplier

import (
	"errors"
)

type Algorithm uint8

const (
	MD5    Algorithm = iota
	SHA1   Algorithm = iota
	SHA224 Algorithm = iota
	SHA256 Algorithm = iota
	SHA384 Algorithm = iota
	SHA512 Algorithm = iota
)

var Algorithms = []Algorithm{
	MD5, SHA1, SHA224, SHA256, SHA384, SHA512,
}

var algoStr = []string{
	"MD5",
	"SHA1",
	"SHA224",
	"SHA256",
	"SHA384",
	"SHA512",
}

func (a Algorithm) Name() string {
	return algoStr[a]
}

func (a Algorithm) Ordinal() uint8 {
	return uint8(a)
}

func ValueOf(val *string) (Algorithm, error) {
	for i, s := range algoStr {
		if s == *val {
			return Algorithm(i), nil
		}
	}

	return 127, errors.New("Algorithm not found")
}
