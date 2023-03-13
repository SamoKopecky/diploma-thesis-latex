type SignAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Verify(pk, msg, signature []byte) bool
	Sign(sk, msg []byte) (signature []byte)
	SignLen() (signLen int)
	PkLen() (pkLen int)
	SkLen() (skLen int)
	Id() (id uint8)
}