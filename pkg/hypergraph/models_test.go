package hypergraph

import "testing"

func BenchmarkUniformERGraph(b *testing.B) {
	n := 1000
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		UniformERGraph(n, 0, 2)
	}
}