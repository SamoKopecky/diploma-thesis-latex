type SignAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Verify(pk, msg, signature []byte) bool
	Sign(sk, msg []byte) (signature []byte)
	SignLen() (signLen int)
	PuKLen() (pkLen int)
	PrKLen() (skLen int)
	Id() (id uint8)
}
