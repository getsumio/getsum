package supplier

type Algorithm uint8

const (
	MD4 Algorithm = iota
	MD5
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	RIPEMD160
	SHA3_224
	SHA3_256
	SHA3_384
	SHA3_512
	SHA512_224
	SHA512_256
	BLAKE2s_256
	BLAKE2b_256
	BLAKE2b_384
	BLAKE2b_512
)

var Algorithms = []Algorithm{
	MD4,
	MD5,
	SHA1,
	SHA224,
	SHA256,
	SHA384,
	SHA512,
	RIPEMD160,
	SHA3_224,
	SHA3_256,
	SHA3_384,
	SHA3_512,
	SHA512_224,
	SHA512_256,
	BLAKE2s_256,
	BLAKE2b_256,
	BLAKE2b_384,
	BLAKE2b_512,
}

var algoStr = []string{
	"MD4",
	"MD5",
	"SHA1",
	"SHA224",
	"SHA256",
	"SHA384",
	"SHA512",
	"RIPEMD160",
	"SHA3_224",
	"SHA3_256",
	"SHA3_384",
	"SHA3_512",
	"SHA512_224",
	"SHA512_256",
	"BLAKE2s_256",
	"BLAKE2b_256",
	"BLAKE2b_384",
	"BLAKE2b_512",
}

func (a Algorithm) Name() string {
	return algoStr[a]
}

func (a Algorithm) Ordinal() uint8 {
	return uint8(a)
}

func ValueOf(val *string) Algorithm {
	for i, s := range algoStr {
		if s == *val {
			return Algorithm(i)
		}
	}

	return 127
}
