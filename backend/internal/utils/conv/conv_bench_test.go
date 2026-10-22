package conv

import "testing"

func BenchmarkStringToInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = StringToInt("12345")
	}
}

func BenchmarkAbsInt(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = AbsInt(-12345)
	}
}
