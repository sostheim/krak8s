/*
Copyright 2017 Samsung SDSA CNCT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"math/rand"
	"time"
)

const hexBytes = "abcdef0123456789"
const (
	// 4 bits to represent a hex value index 0 - 16
	hexIdxBits = 4
	// all 1-bits, as many as hexIdxBits
	hexIdxMask = 1<<hexIdxBits - 1
	// # of hex indices fitting in 16 bits
	hexIdxMax = 16 / hexIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
