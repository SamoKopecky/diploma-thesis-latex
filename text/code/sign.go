type SignAlgorithm interface {
	KeyGen() (puK, prK []byte)
	Verify(puK, msg, signature []byte) bool
	Sign(prK, msg []byte) (signature []byte)
	SignLen() (signLen int)
	PuKLen() (pkLen int)
	PrKLen() (skLen int)
	Id() (id uint8)
}
