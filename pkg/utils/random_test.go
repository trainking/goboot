package utils

import "testing"

func TestRandStringNumber(t *testing.T) {
	s, _ := RandStringNumber(6)

	t.Log(s)
}

func BenchmarkRandStringNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringNumber(6)
	}
}

func BenchmarkRandStringHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStringHex(6)
	}
}
