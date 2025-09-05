//go:build demo
// +build demo

// Package framework contains tests that verify the framework's behavior
// when properties fail intentionally. These tests ensure that the framework
// correctly handles failures, shrinking, and parallel execution paths.
package framework

import (
	"math/rand"
	"testing"

	"arcsyn.io/propx"
)

// TestForAll_ShrinkingFailure tests the shrinking mechanism with intentional failures.
// This test verifies that the framework correctly shrinks values when properties fail.
func TestForAll_ShrinkingFailure(t *testing.T) {
	config := propx.Config{
		Seed:        12345,
		Examples:    1,
		MaxShrink:   2,
		ShrinkStrat: "bfs",
		Parallelism: 1,
	}

	gen := propx.From(func(r *rand.Rand, sz propx.Size) (int, propx.Shrinker[int]) {
		return 42, func(accept bool) (int, bool) {
			return 0, false
		}
	})

	propx.ForAll(t, config, gen)(func(t *testing.T, val int) {
		t.Errorf("This should fail: got %d", val)
	})
}
