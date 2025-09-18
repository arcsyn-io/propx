// File: gen/pair.go
package gen

import (
	"math/rand"
)

// Pair represents a pair of values of types A and B.
type Pair[A, B any] struct {
	First  A
	Second B
}

// PairOf creates a generator that produces pairs of values from two generators.
// The generated pairs will have shrinking capabilities that try to shrink both
// components independently.
//
// Example usage:
//
//	// Generate pairs of integers
//	pairGen := gen.PairOf(gen.Int(), gen.Int())
//	
//	// Use in property-based testing
//	propx.ForAll(t, cfg, pairGen)(func(t *testing.T, p gen.Pair[int, int]) {
//		// Test property using p.First and p.Second
//		if p.First+p.Second != p.Second+p.First {
//			t.Errorf("addition is not commutative for (%d, %d)", p.First, p.Second)
//		}
//	})
func PairOf[A, B any](ga Generator[A], gb Generator[B]) Generator[Pair[A, B]] {
	return From(func(r *rand.Rand, sz Size) (Pair[A, B], Shrinker[Pair[A, B]]) {
		if r == nil {
			r = rand.New(rand.NewSource(rand.Int63())) // #nosec G404 -- Using math/rand for deterministic property-based testing
		}
		
		// Generate both values
		a, sa := ga.Generate(r, sz)
		b, sb := gb.Generate(r, sz)
		
		pair := Pair[A, B]{First: a, Second: b}
		
		// Shrinker strategy: try shrinking both components
		// First try shrinking the first component, then the second
		shrinkingFirst := true
		currentA := a
		currentB := b
		
		return pair, func(accept bool) (Pair[A, B], bool) {
			if shrinkingFirst {
				// Try shrinking the first component
				if na, ok := sa(accept); ok {
					currentA = na
					return Pair[A, B]{First: na, Second: currentB}, true
				}
				// First component exhausted, switch to second
				shrinkingFirst = false
				accept = false // Reset accept for second component
			}
			
			// Try shrinking the second component
			if nb, ok := sb(accept); ok {
				currentB = nb
				return Pair[A, B]{First: currentA, Second: nb}, true
			}
			
			// Both components exhausted
			var zero Pair[A, B]
			return zero, false
		}
	})
}

// Tuple is an alias for Pair for better readability in some contexts.
type Tuple[A, B any] = Pair[A, B]

// TupleOf is an alias for PairOf for better readability in some contexts.
func TupleOf[A, B any](ga Generator[A], gb Generator[B]) Generator[Tuple[A, B]] {
	return PairOf(ga, gb)
}