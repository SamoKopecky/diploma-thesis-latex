package dilithium#

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	"golang.org/x/crypto/sha3"
)

func modP(a, mod int) (o int) {
	o = ((a % mod) + mod) % mod
	return
}

func modPM(a, mod int) int {
	first := int(math.Abs(float64(((a % mod) - mod) % mod)))
	if first <= mod/2 {
		return int(-first)
	}
	return int(modP(a, mod))
}

func powerToRound(r, d int) (int, int) {
	r = modP(r, Q)
	r0 := modPM(r, 1<<D)
	return (r - r0) / (1 << D), r0
}

func decompose(r, alpha int) (r_1, r_0 int) {
	r = modP(r, Q)
	r_0 = modPM(r, alpha)
	if r-r_0 == Q-1 {
		r_1 = 0
		r_0 -= 1
	} else {
		r_1 = int((r - r_0) / alpha)
	}
	return
}

func highBits(r, alpha int) (r_1 int) {
	r_1, _ = decompose(r, alpha)
	return
}

func lowBits(r, alpha int) (r_0 int) {
	_, r_0 = decompose(r, alpha)
	return
}

func makeHint(z, r, alpha int) bool {
	return highBits(r, alpha) != highBits(r+z, alpha)
}

func useHint(h bool, r, alpha int) int {
	m := int((Q - 1) / alpha)
	r1, r0 := decompose(r, alpha)
	if h && r0 > 0 {
		return modP(r1+1, m)
	} else if h && r0 <= 0 {
		return modP(r1-1, m)
	}
	return r1
}

func sampleInBall(c_wave []byte) (c []int) {
	c = make([]int, 256)
	shake := sha3.NewShake256()
	shake.Write(c_wave)
	o := make([]byte, 8)
	shake.Read(o)

	bits := bytesToBits(o)[:Tau]
	for i := 256 - Tau; i < 256; i++ {
		j_byte := make([]byte, 1)
		j := byte(255)
		for j > byte(i) {
			shake.Read(j_byte)
			j = j_byte[0]
		}
		s := bits[i-(256-Tau)]
		c[i] = c[j]
		c[j] = int(1 - int8(2*s))
	}
	return
}

func bytesToBits(bytes []byte) (bits []byte) {
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			bits = append(bits, extractBit(int(bytes[i]), j))
		}
	}
	return
}

func polyToBits(poly []int, l int) (bits []byte) {
	for i := 0; i < 256; i++ {
		for j := 0; j < l; j++ {
			bits = append(bits, extractBit((poly[i]), j))
		}
	}
	return
}

func extractBit(from int, power int) (bit byte) {
	bit = byte(from & (1 << power) >> power)
	return
}

func genRand(bits int) (xofOutput []byte) {
	len := bits / 8
	xofOutput = make([]byte, len)
	randBytes := make([]byte, 256)
	rand.Read(randBytes)
	xofOutput = shake256(randBytes, len)
	return
}

func expandS(ro_dash []byte) (vectors [][]int) {
	for i := 0; i < L+K; i++ {
		poly := []int{}
		shake := sha3.NewShake256()
		i_bytes := make([]byte, 2)
		i_bytes[0] = byte(i)
		i_bytes[1] = byte(i >> 8)
		shake.Write(ro_dash)
		shake.Write(i_bytes)

		for len(poly) < N {
			o := make([]byte, 1)
			shake.Read(o)

			two_ints := [...]uint8{uint8(o[0]) >> 4, uint8(o[0]) & 0xF}

			if Eta == 2 {
				for _, v := range two_ints {
					if v >= 15 {
						continue
					}
					if v < 0 {
						continue
					}
					poly = append(poly, Eta-(int(v)%5))
				}
			}
			// TODO: eta = 4
		}
		vectors = append(vectors, poly[:N])
	}
	return
}

func expandMask(ro_dash []byte, kappa int) (y [][]int) {
	for i := 0; i < L; i++ {
		poly := make([]int, N)
		shake := sha3.NewShake256()
		shake.Write(ro_dash)
		sum := kappa + i
		bytes := make([]byte, 2)
		bytes[0] = byte(sum)
		bytes[1] = byte(sum >> 8)
		shake.Write(bytes)

		for j := 0; j < N; j++ {
			o := make([]byte, 4)
			shake.Read(o)
			o[0] = 0
			o[1] = o[1] & 2
			poly[j] = GammaOne - int(binary.BigEndian.Uint32(o))
		}
		y = append(y, poly)
	}
	return
}

func BytesEqual(a, b []byte) (equal bool) {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
