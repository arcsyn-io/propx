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

// TestForAll_ParallelFailure tests failure scenarios in runParallel.
// This test verifies that the framework correctly handles failures in parallel mode.
func TestForAll_ParallelFailure(t *testing.T) {
	config := propx.Config{
		Seed:        12345,
		Examples:    3,
		MaxShrink:   5,
		ShrinkStrat: "bfs",
		Parallelism: 2,
	}

	gen := propx.From(func(r *rand.Rand, sz propx.Size) (int, propx.Shrinker[int]) {
		return 42, func(accept bool) (int, bool) {
			return 0, false
		}
	})

	// This should fail and trigger the parallel failure path
	propx.ForAll(t, config, gen)(func(t *testing.T, val int) {
		t.Errorf("This should fail: got %d", val)
	})
}

// TestForAll_ParallelFailureWithShrinking tests parallel failure with shrinking.
// This test verifies that the framework correctly handles shrinking in parallel mode.
func TestForAll_ParallelFailureWithShrinking(t *testing.T) {
	config := propx.Config{
		Seed:        12345,
		Examples:    2,
		MaxShrink:   3,
		ShrinkStrat: "bfs",
		Parallelism: 2,
	}

	shrinkerCallCount := 0
	gen := propx.From(func(r *rand.Rand, sz propx.Size) (int, propx.Shrinker[int]) {
		return 5, func(accept bool) (int, bool) {
			shrinkerCallCount++
			if shrinkerCallCount <= 2 {
				return shrinkerCallCount, true
			}
			return 0, false
		}
	})

	// This should fail and trigger parallel shrinking
	propx.ForAll(t, config, gen)(func(t *testing.T, val int) {
		t.Errorf("This should fail: got %d", val)
	})
}

// TestForAll_ParallelStopOnFirstFailureFalse tests parallel execution
// with StopOnFirstFailure set to false.
func TestForAll_ParallelStopOnFirstFailureFalse(t *testing.T) {
	config := propx.Config{
		Seed:               12345,
		Examples:           3,
		MaxShrink:          2,
		ShrinkStrat:        "bfs",
		Parallelism:        2,
		StopOnFirstFailure: false,
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
