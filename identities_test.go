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
	"testing"
)

func TestRandHexString(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	// Generate N Unique Id's without a duplicate
	iterations := 100000000
	check := make(map[string]bool, iterations)
	dups := make(map[string]bool)
	for i := 0; i < iterations; i++ {
		tmp := RandHexString(12)
		if _, found := check[tmp]; !found {
			check[tmp] = true
		} else if _, dup := dups[tmp]; !dup {
			dups[tmp] = true
		} else {
			t.Errorf("RandHexString(), two strike duplicate (%d, c[%t], d[%t], %s), want unique",
				i, check[tmp], dups[tmp], tmp)
		}
	}
	percent := (float64(len(dups)) * float64(100)) / float64(iterations)
	if percent > 0.05 {
		t.Errorf("RandHexString(), excessive rate of duplicates (%f), want < 0.05%%", percent)
	}
}

func BenchmarkRandHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandHexString(8)
	}
}
