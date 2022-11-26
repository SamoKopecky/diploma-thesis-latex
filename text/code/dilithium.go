package dilithium#

func KeyGen() (pk, sk []byte) {
	zeta := shake256(genRand(256), 128)
	rho := zeta[:32]
	rho_dash := zeta[32:96]
	K := zeta[96:]

	A_hat := expandA(rho)
	s_1, s_2 := expandS(rho_dash)
	for i := 0; i < L; i++ {
		t[i] = mulPolyVec(A_hat[i], nttPolyVec(s_1))
	}
	t = addPolyVec(invNttPolyVec(t), s_2)
	t_1, t_0 := powerToRoundPolyVec(t, D)

	pk = append(rho, bitPack(t_1, 10))
	sk = append(rho, K)
	sk = append(sk, shake256(pk, 32))
	sk = append(sk, bitPack(s_1, 3))
	sk = append(sk, bitPack(s_2, 3))
	sk = append(sk, bitPack(t_0, 13))
	return
}

func Sign(sk []byte, message []byte) (sigma []byte) {
	rho := sk[:32]
	K := sk[32:64]
	tr := sk[64:96]
	s_1 := bitUnpack(sk[96:SBytes+96], 3)
	s_2 := bitUnpack(sk[96+SBytes:SBytes*2+96], 3)
	t_0 := bitUnpack(sk[96+SBytes*2:], 13)
	A_hat := expandA(rho)
	s_1_hat := nttPolyVec(s_1)
	s_2_hat := nttPolyVec(s_2)
	t_0_hat := nttPolyVec(t_0)

	shake := append(tr, message...)
	mi := shake256(shake, 64)
	shake = append(K, mi...)
	rho_dash := shake256(shake, 64)

	kappa := -L

	for true {
	kappa += L

	y := expandMask(rho_dash, kappa)
	y_hat := nttPolyVec(y)
	for i := 0; i < L; i++ {
		w[i] = invNtt(mulPolyVec(A_hat[i], y_hat))
	}
	w_1 := highBitsPolyVec(w, 2*GammaTwo)

	shake = append(mi, bitPack(w_1, 6))
	c_wave := shake256(shake, 32)
	c := sampleInBall(c_wave)

	cs_2 := invNttPolyVec(scalePolyVecByPoly(s_2_hat, c_hat))
	cs_1 := invNttPolyVec(scalePolyVecByPoly(s_1_hat, c_hat))
	z := addPolyVec(y, cs_1),

	w_sub_cs_2 := subPolyVec(w, cs_2)
	r_0 := lowBitsPolyVec(w_sub_cs_2, 2*GammaTwo)
	if checkNormPolyVec(z, GammaOne-Beta) || checkNormPolyVec(r_0, GammaTwo-Beta) {
		continue
	}

	ct_0 := invNttPolyVec(scalePolyVecByPoly(t_0_hat, c_hat))
	ct_0_inv := inversePolyVec(ct_0)
	added_vecs := addPolyVec(w_sub_cs_2, ct_0)
	h := makeHintPolyVec(ct_0_inv, added_vecs, GammaTwo*2)

	ones := 0
	for i := 0; i < len(h); i++ {
		for j := 0; j < N; j++ {
			if h[i][j] == 1 {
				ones++
			}
		}
	}
	if checkNormPolyVec(modPMPolyVec(ct_0, Q), GammaTwo) {
		continue
	}
	if ones > Omega {
		continue
	}

	z_packed := bitPackAlteredPolyVec(z, GammaOne, 18)
	h_packed := bitPackHint(h)
	sigma = append(c_wave, z_packed...)
	sigma = append(sigma, h_packed...)
	break
	}
	return
}

func Verify(pk, message, sigma []byte) (verified bool) {
	rho := pk[:32]
	t_1_bytes := pk[32:]
	c_wave := sigma[:32]
	z := bitUnpackAlteredPolyVec(sigma[32:ZBytes+32], GammaOne, 18)
	h := bitUnpackHint(sigma[ZBytes+32:])
	t_1 := bitUnpackPolyVec(t_1_bytes, 10)

	A_hat := expandA(rho)

	shake := append(rho, t_1_bytes...)
	shake = append(shake256(shake, 32), message...)
	mi := shake256(shake, 64)

	c := sampleInBall(c_wave)
	c_hat := ntt(c)

	z_hat := nttPolyVec(z)

	for i := 0; i < L; i++ {
		Az[i] = mulPolyVec(A_hat[i], z_hat)
	}
	t_1_scaled_hat := nttPolyVec(scalePolyVecByInt(t_1, 1<<D))
	ct_1 := scalePolyVecByPoly(t_1_scaled_hat, c_hat)

	r := invNttPolyVec(subPolyVec(Az, ct_1))
	w_1 := useHintPolyVec(h, r, 2*GammaTwo)
	shake = append(mi, bitUnpack(w_1, 6))
	verified = BytesEqual(c, shake256(shake, 32))
	return
}

func makeHint(z, r, alpha int) bool {
	return highBits(r, alpha) != highBits(r+z, alpha)
}
