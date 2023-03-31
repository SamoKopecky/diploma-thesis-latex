type KemAlgorithm interface {
	KeyGen() (puK, prK []byte)
	Dec(c, prK []byte) (key []byte)
	Enc(puK []byte) (c, key []byte)
	EkLen() (ekLen int)
	CLen() (cLen int)
	Id() (id uint8)
}

var kems = map[string]KemAlgorithm{
	"PqComKyber512": &kem.PqComKyber512{},
	"CirclKyber512": &kem.CirclKyber512{},
}
