package fib

import "testing"

// START OMIT
func Fib(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		a, b = b, a+b
	}
	return a
}

var Result int

func BenchmarkFibWrong(b *testing.B) {
	Result = Fib(b.N)
}

func BenchmarkFibWrong2(b *testing.B) {
	var r int
	for n := 0; n < b.N; n++ {
		r = Fib(n)
	}
	Result = r
}

// END OMIT
