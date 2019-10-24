package supplier

type Algorithm uint8

const (
	MD2 Algorithm = iota
	MD4
	MD5
	GOST
	SHA1
	SHA224
	SHA256
	SHA384
	SHA512
	RMD160
	SHA3_224
	SHA3_256
	SHA3_384
	SHA3_512
	SHA512_224
	SHA512_256
	BLAKE2s256
	BLAKE2b256
	BLAKE2b384
	BLAKE2b512
	SHAKE128
	SHAKE256
	SM3
)

var Algorithms = []Algorithm{
	MD2,
	MD4,
	MD5,
	GOST,
	SHA1,
	SHA224,
	SHA256,
	SHA384,
	SHA512,
	RMD160,
	SHA3_224,
	SHA3_256,
	SHA3_384,
	SHA3_512,
	SHA512_224,
	SHA512_256,
	BLAKE2s256,
	BLAKE2b256,
	BLAKE2b384,
	BLAKE2b512,
	SHAKE128,
	SHAKE256,
	SM3,
}

var algoStr = []string{
	"MD2",
	"MD4",
	"MD5",
	"GOST",
	"SHA1",
	"SHA224",
	"SHA256",
	"SHA384",
	"SHA512",
	"RMD160",
	"SHA3-224",
	"SHA3-256",
	"SHA3-384",
	"SHA3-512",
	"SHA512-224",
	"SHA512-256",
	"BLAKE2s256",
	"BLAKE2b256",
	"BLAKE2b384",
	"BLAKE2b512",
	"SHAKE128",
	"SHAKE256",
	"SM3",
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
