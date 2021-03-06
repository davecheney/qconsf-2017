// +build ignore

package main

import "testing"

const m1 = 0x5555555555555555
const m2 = 0x3333333333333333
const m4 = 0x0f0f0f0f0f0f0f0f
const h01 = 0x0101010101010101

func popcnt(x uint64) int {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return int((x * h01) >> 56)
}

// START OMIT
func BenchmarkPopcnt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// optimised away
	}
}

// END OMIT
