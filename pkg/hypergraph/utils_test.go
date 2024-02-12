package hypergraph

import (
	"log"
	"testing"
)

func TestBinomialCoefficient(t *testing.T) {
	c := binomialCoefficient(1000, 3)
	exp := 166167000

	if c != exp {
		log.Fatalf("Wrong solution, got %d, expected %d.", c, exp)
	}
}