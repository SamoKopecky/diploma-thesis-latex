package kyber

func cpapkeKeyGen() (pk []byte, sk []byte) {
	rho, sigma := randBytes(64)
	A_hat := genPolyMat(rho)
	s_hat := nttPolyVec(randPolyVec(sigma))
	e_hat := nttPolyVec(randPolyVec(sigma))
	for i := 0; i < K; i++ {
		t_hat[i] = mulPolyVec(A_hat[i], s_hat)
	}
	t_hat = addPolyVec(t_hat, e_hat)
	sk = encodePolyVec(s_hat)
	pk = encodePolyVec(t_hat)
	pk = append(pk, rho...)
	return
}

func cpapkeEnc(pk []byte, m []byte, r []byte) (c []byte) {
	t_hat := decodePolyVec(pk)
	// Get the last 32 bytes from PK
	rho := pk[len(pk)-32:]
	A_hat := genPolyMat(rho)
	r_hat := nttPolyVec(randPolyVec(r))
	e1 := randPolyVec(r)
	e2 := randPoly(r)
	for i := 0; i < K; i++ {
		u[i] = mulPolyVec(A_hat[i], r_hat)
	}
	u = addPolyVec(invNttPolyVec(u), e1)
	parsed_m := decompress(decode(m))
	v := invNtt(mulPolyVec(t_hat, r_hat))
	v = polyAdd(v, e2)
	v = polyAdd(v, parsed_m)
	for i := 0; i < K; i++ {
		c1 = append(c1, encode(compress(u[i], Du), Du))
	}
	c2 := encode(compress(v, Dv), Dv)
	c = append(c1, c2)
	return
}

func cpapkeDec(sk []byte, c []byte) (m []byte) {

	c2 := c[Du * K * N / 8:]
	u_decoded := decodePolyVec(c, Du)

	for i := 0; i < K; i++ {
		u_hat[i] = decompress(u_decoded[i], Du)
	}

	nttPolyVec(u_hat)

	v := decompress(decode(c2, Dv), Dv)
	s_hat := decodePolyVec(sk, 12)

	s_hat_u_hat := mulPolyVec(s_hat, u_hat)
	invNtt(s_hat_u_hat)
	first_m := polySub(v, s_hat_u_hat)

	m = encode(compress(first_m, 1), 1)
	return
}

func CcakemKeyGen() (pk []byte, sk []byte) {
	pk, sk := cpapkeKeyGen()
	sk = append(sk, pk)
	sk = append(sk, hash32(pk)
	sk = append(sk, randBytes(32))
	return
}

func CcakemEnc(pk []byte) (c, key []byte) {
	m := hash32(randBytes(32))

	g_input := []byte{}
	g_input = append(g_input, m...)
	g_input = append(g_input, hash32(pk)...)

	K_dash, r := hash64(g_input)
	c = cpapkeEnc(pk, m, r)

	kdf_input := []byte{}
	kdf_input = append(kdf_input, K_dash...)
	kdf_input = append(kdf_input, hash32(c)...)
	key = kdf(kdf_input, 32)
	return
}

func CcakemDec(c, sk []byte) (key []byte) {
	keySize := 12 * K * N / 8
	pk := sk[keySize : keySize*2+32]
	hash := sk[keySize*2+32 : keySize*2+64]
	z := sk[keySize*2+64:]

	m_dash := cpapkeDec(sk, c)

	g_input := []byte{}
	g_input = append(g_input, m_dash...)
	g_input = append(g_input, hash...)
	k_dash, r_dash := hash64(g_input)

	c_dash := cpapkeEnc(pk, m_dash, r_dash)
	hash_c := hash32(c)

	kdf_input := []byte{}
	if BytesEqual(c, c_dash) {
		kdf_input = append(kdf_input, k_dash...)
		kdf_input = append(kdf_input, hash_c...)
		key = kdf(kdf_input, SharedKeySize)
	} else {
		kdf_input = append(z, k_dash...)
		kdf_input = append(kdf_input, hash_c...)
		key = kdf(kdf_input, SharedKeySize)
	}
	return
}