type KemAlgorithm interface {
	KeyGen() (pk, sk []byte)
	Dec(c, sk []byte) (key []byte)
	Enc(pk []byte) (c, key []byte)
	EkLen() (ekLen int)
	CLen() (cLen int)
	Id() (id uint8)
}

var kems = map[string]KemAlgorithm{
	"PqComKyber512": &kem.PqComKyber512{},
	"CirclKyber512": &kem.CirclKyber512{},
}
