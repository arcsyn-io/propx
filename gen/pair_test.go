// File: gen/pair_test.go
package gen

import (
	"math/rand"
	"testing"
)

func TestPairOf_Generation(t *testing.T) {
	// Create generators for int and string
	intGen := IntRange(1, 100) // Use IntRange for exact bounds
	strGen := StringAlpha(Size{Min: 1, Max: 10})
	
	// Create pair generator
	pairGen := PairOf(intGen, strGen)
	
	// Generate some pairs
	r := rand.New(rand.NewSource(42))
	for i := 0; i < 10; i++ {
		pair, shrinker := pairGen.Generate(r, Size{})
		
		// Check that the pair has valid values
		if pair.First < 1 || pair.First > 100 {
			t.Errorf("First component out of range: %d", pair.First)
		}
		if len(pair.Second) < 1 || len(pair.Second) > 10 {
			t.Errorf("Second component length out of range: %d", len(pair.Second))
		}
		
		// Check that shrinker is not nil
		if shrinker == nil {
			t.Error("Shrinker should not be nil")
		}
	}
}

func TestPairOf_Shrinking(t *testing.T) {
	// Create a simple generator that generates values that can be shrunk
	intGen := From(func(r *rand.Rand, sz Size) (int, Shrinker[int]) {
		val := 10 // Start with 10
		return val, func(accept bool) (int, bool) {
			if val > 0 {
				val--
				return val, true
			}
			return 0, false
		}
	})
	
	strGen := From(func(r *rand.Rand, sz Size) (string, Shrinker[string]) {
		val := "hello" // Start with "hello"
		return val, func(accept bool) (string, bool) {
			if len(val) > 1 {
				val = val[:len(val)-1] // Remove last character
				return val, true
			}
			return "", false
		}
	})
	
	pairGen := PairOf(intGen, strGen)
	r := rand.New(rand.NewSource(42))
	
	pair, shrinker := pairGen.Generate(r, Size{})
	
	// Test shrinking
	originalFirst := pair.First
	originalSecond := pair.Second
	
	// First shrink should affect the first component
	shrunkPair, ok := shrinker(true)
	if !ok {
		t.Error("Expected shrinking to succeed")
	}
	if shrunkPair.First >= originalFirst {
		t.Errorf("Expected first component to shrink from %d, got %d", originalFirst, shrunkPair.First)
	}
	if shrunkPair.Second != originalSecond {
		t.Errorf("Expected second component to remain %q, got %q", originalSecond, shrunkPair.Second)
	}
	
	// Continue shrinking first component until exhausted
	for i := 0; i < 15; i++ {
		shrunkPair, ok = shrinker(true)
		if !ok {
			break
		}
		if shrunkPair.First >= originalFirst {
			t.Errorf("Expected first component to continue shrinking, got %d", shrunkPair.First)
		}
	}
	
	// Now try to shrink again - should start working on second component
	shrunkPair, ok = shrinker(true)
	if ok {
		// If we get here, it should be shrinking the second component
		if shrunkPair.Second == originalSecond {
			t.Errorf("Expected second component to start shrinking from %q", originalSecond)
		}
	}
}

func TestPairOf_EmptyShrinking(t *testing.T) {
	// Create generators that don't shrink
	constGen := Const(42)
	constStrGen := Const("test")
	
	pairGen := PairOf(constGen, constStrGen)
	r := rand.New(rand.NewSource(42))
	
	_, shrinker := pairGen.Generate(r, Size{})
	
	// Shrinking should fail immediately
	_, ok := shrinker(true)
	if ok {
		t.Error("Expected shrinking to fail for constant generators")
	}
}

func TestTupleOf_Alias(t *testing.T) {
	// Test that TupleOf is equivalent to PairOf
	intGen := Int(Size{Min: 1, Max: 10})
	strGen := StringAlpha(Size{Min: 1, Max: 5})
	
	pairGen := PairOf(intGen, strGen)
	tupleGen := TupleOf(intGen, strGen)
	
	r1 := rand.New(rand.NewSource(42))
	r2 := rand.New(rand.NewSource(42))
	
	pair, _ := pairGen.Generate(r1, Size{})
	tuple, _ := tupleGen.Generate(r2, Size{})
	
	// They should generate equivalent values
	if pair.First != tuple.First {
		t.Errorf("Pair and Tuple should generate same first component: %d vs %d", pair.First, tuple.First)
	}
	if pair.Second != tuple.Second {
		t.Errorf("Pair and Tuple should generate same second component: %q vs %q", pair.Second, tuple.Second)
	}
}

func TestPair_StructAccess(t *testing.T) {
	// Test that we can access the fields of the Pair struct
	intGen := Int(Size{Min: 1, Max: 10})
	strGen := StringAlpha(Size{Min: 1, Max: 5})
	
	pairGen := PairOf(intGen, strGen)
	r := rand.New(rand.NewSource(42))
	
	pair, _ := pairGen.Generate(r, Size{})
	
	// Test field access
	if pair.First == 0 {
		t.Error("First field should not be zero")
	}
	if pair.Second == "" {
		t.Error("Second field should not be empty")
	}
}