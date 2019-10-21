package supplier

type Algorithm int

const (
	MD5 = iota
	SHA1
	SHA256
	SHA512
	SHA3
	BLAKE2s
	BLAKE2b
)
