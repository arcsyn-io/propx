// File: examples/pair_test.go
package examples

import (
	"testing"

	"arcsyn.io/propx"
)

// TestPairAdditionCommutativity demonstrates how to use the Pair generator
// to test properties that involve two parameters.
func TestPairAdditionCommutativity(t *testing.T) {
	// Create a generator for pairs of integers
	pairGen := propx.PairOf(
		propx.Int(propx.Size{Min: -100, Max: 100}),
		propx.Int(propx.Size{Min: -100, Max: 100}),
	)

	// Test the commutative property of addition
	propx.ForAll(t, propx.Default(), pairGen)(func(t *testing.T, p propx.Pair[int, int]) {
		// Property: a + b == b + a (commutativity)
		if p.First+p.Second != p.Second+p.First {
			t.Errorf("addition is not commutative for (%d, %d): %d + %d = %d, but %d + %d = %d",
				p.First, p.Second,
				p.First, p.Second, p.First+p.Second,
				p.Second, p.First, p.Second+p.First)
		}
	})
}

// TestPairMultiplicationAssociativity demonstrates testing associativity
// with pairs of integers.
func TestPairMultiplicationAssociativity(t *testing.T) {
	// Generate pairs of small integers to avoid overflow
	pairGen := propx.PairOf(
		propx.Int(propx.Size{Min: -10, Max: 10}),
		propx.Int(propx.Size{Min: -10, Max: 10}),
	)

	propx.ForAll(t, propx.Default(), pairGen)(func(t *testing.T, p propx.Pair[int, int]) {
		// Property: (a * b) * c == a * (b * c) for some constant c
		c := 2
		left := (p.First * p.Second) * c
		right := p.First * (p.Second * c)

		if left != right {
			t.Errorf("multiplication is not associative for (%d, %d) with c=%d: (%d * %d) * %d = %d, but %d * (%d * %d) = %d",
				p.First, p.Second, c,
				p.First, p.Second, c, left,
				p.First, p.Second, c, right)
		}
	})
}

// TestPairStringConcatenation demonstrates using pairs of strings.
func TestPairStringConcatenation(t *testing.T) {
	// Generate pairs of strings
	pairGen := propx.PairOf(
		propx.StringAlpha(propx.Size{Min: 1, Max: 10}),
		propx.StringAlpha(propx.Size{Min: 1, Max: 10}),
	)

	propx.ForAll(t, propx.Default(), pairGen)(func(t *testing.T, p propx.Pair[string, string]) {
		// Property: concatenation length is sum of individual lengths
		concatenated := p.First + p.Second
		expectedLength := len(p.First) + len(p.Second)

		if len(concatenated) != expectedLength {
			t.Errorf("concatenation length mismatch for (%q, %q): got length %d, expected %d",
				p.First, p.Second, len(concatenated), expectedLength)
		}
	})
}

// TestPairMixedTypes demonstrates using pairs of different types.
func TestPairMixedTypes(t *testing.T) {
	// Generate pairs of int and string
	pairGen := propx.PairOf(
		propx.Int(propx.Size{Min: 1, Max: 100}),
		propx.StringDigits(propx.Size{Min: 1, Max: 5}),
	)

	propx.ForAll(t, propx.Default(), pairGen)(func(t *testing.T, p propx.Pair[int, string]) {
		// Property: if we convert the string to int and add to the first int,
		// the result should be a valid integer
		if len(p.Second) > 0 {
			// This is a simple property - we're just checking that
			// the string contains only digits (which it should, given our generator)
			for _, char := range p.Second {
				if char < '0' || char > '9' {
					t.Errorf("string %q contains non-digit character %c", p.Second, char)
				}
			}
		}
	})
}

// TestTupleAlias demonstrates using the Tuple alias for better readability.
func TestTupleAlias(t *testing.T) {
	// Use TupleOf instead of PairOf for better readability
	tupleGen := propx.TupleOf(
		propx.Bool(),
		propx.Bool(),
	)

	propx.ForAll(t, propx.Default(), tupleGen)(func(t *testing.T, tup propx.Tuple[bool, bool]) {
		// Property: logical AND is commutative
		if (tup.First && tup.Second) != (tup.Second && tup.First) {
			t.Errorf("logical AND is not commutative for (%t, %t)",
				tup.First, tup.Second)
		}
	})
}
