package main

import (
	"math/rand"
	"time"
)

// Implementations
func init() {
	rand.Seed(time.Now().UnixNano())
}

const hexBytes = "abcdef0123456789"
const (
	hexIdxBits = 4                 // 4 bits to represent a hex value index
	hexIdxMask = 1<<hexIdxBits - 1 // All 1-bits, as many as hexIdxBits
	hexIdxMax  = 16 / hexIdxBits   // # of hex indices fitting in 16 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandHexString - generate a random hex string of n characters.
func RandHexString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits as an int64
	for i, cache, remain := n-1, src.Int63(), hexIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), hexIdxMax
		}
		if idx := int(cache & hexIdxMask); idx < len(hexBytes) {
			b[i] = hexBytes[idx]
			i--
		}
		cache >>= hexIdxBits
		remain--
	}

	return string(b)
}
