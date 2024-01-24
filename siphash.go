package siphash128

import "encoding/binary"

// SipHash128ch: sipHash128() in ClickHouse, their own implementation of SipHash:
// https://clickhouse.com/codebrowser/ClickHouse/src/Common/SipHash.h.html

// Func based on:
// https://github.com/dchest/siphash/blob/master/hash.go

func SipHash128(b []byte) [16]byte {
	var k0, k1 uint64
	var BlockSize = 8

	v0 := k0 ^ 0x736f6d6570736575
	v1 := k1 ^ 0x646f72616e646f6d
	v2 := k0 ^ 0x6c7967656e657261
	v3 := k1 ^ 0x7465646279746573
	t := uint64(len(b)) << 56

	// compress
	for len(b) >= BlockSize {
		m := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56

		v3 ^= m
		v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
		v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
		v0 ^= m

		b = b[BlockSize:]
	}

	// compress last block
	switch len(b) {
	case 7:
		t |= uint64(b[6]) << 48
		fallthrough
	case 6:
		t |= uint64(b[5]) << 40
		fallthrough
	case 5:
		t |= uint64(b[4]) << 32
		fallthrough
	case 4:
		t |= uint64(b[3]) << 24
		fallthrough
	case 3:
		t |= uint64(b[2]) << 16
		fallthrough
	case 2:
		t |= uint64(b[1]) << 8
		fallthrough
	case 1:
		t |= uint64(b[0])
	}

	// finalize
	v3 ^= t
	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
	v0 ^= t

	v2 ^= 0xff

	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)
	v0, v1, v2, v3 = sipHash128chRound(v0, v1, v2, v3)

	// compute hash
	hash := [16]byte{}
	binary.LittleEndian.PutUint64(hash[0:], v0^v1)
	binary.LittleEndian.PutUint64(hash[8:], v2^v3)

	return hash
}

func sipHash128chRound(v0 uint64, v1 uint64, v2 uint64, v3 uint64) (uint64, uint64, uint64, uint64) {
	v0 += v1
	v1 = v1<<13 | v1>>(64-13)
	v1 ^= v0
	v0 = v0<<32 | v0>>(64-32)

	v2 += v3
	v3 = v3<<16 | v3>>(64-16)
	v3 ^= v2

	v0 += v3
	v3 = v3<<21 | v3>>(64-21)
	v3 ^= v0

	v2 += v1
	v1 = v1<<17 | v1>>(64-17)
	v1 ^= v2
	v2 = v2<<32 | v2>>(64-32)

	return v0, v1, v2, v3
}
