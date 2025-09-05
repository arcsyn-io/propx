// Package propx provides property-based testing functionality for Go.
// It allows you to test properties of your code by generating random test cases
// and automatically shrinking counterexamples when failures are found.
//
// This is the main entry point for the propx library. It re-exports the most
// commonly used types and functions from the internal packages to provide a
// clean and simple API for users.
//
// Example usage:
//
//	import "arcsyn.io/propx"
//
//	func TestAdditionIdentity(t *testing.T) {
//		propx.ForAll(t, propx.Default(), propx.Int())(func(t *testing.T, x int) {
//			if x+0 != x {
//				t.Errorf("addition identity failed for %d", x)
//			}
//		})
//	}
package propx

import (
	"math/rand"
	"testing"

	"arcsyn.io/propx/gen"
	"arcsyn.io/propx/gen/domain"
	"arcsyn.io/propx/prop"
	"arcsyn.io/propx/quick"
)

// =============================================================================
// PROPERTY-BASED TESTING
// =============================================================================

// Config holds the configuration for property-based testing.
type Config = prop.Config

// Default returns a default configuration for property-based testing.
// This configuration uses sensible defaults and can be customized via
// command-line flags or by modifying the returned Config struct.
func Default() Config {
	return prop.Default()
}

// ForAll runs a property-based test with the given configuration and generator.
// It generates test cases using the provided generator and runs the property
// function for each generated value. If a counterexample is found, it will
// be automatically shrunk to find a minimal failing case.
//
// Parameters:
//   - t: The testing.T instance for the current test
//   - cfg: Configuration for the property-based test
//   - gen: Generator that produces test values
//   - property: Function that defines the property to test
//
// Example:
//
//	propx.ForAll(t, propx.Default(), propx.Int())(func(t *testing.T, x int) {
//		if x+0 != x {
//			t.Errorf("addition identity failed for %d", x)
//		}
//	})
func ForAll[T any](t *testing.T, cfg Config, g gen.Generator[T]) func(func(*testing.T, T)) {
	return prop.ForAll(t, cfg, g)
}

// =============================================================================
// STATE MACHINE TESTING
// =============================================================================

// StateMachine represents a state machine for property-based testing.
type StateMachine[S, C any] = prop.StateMachine[S, C]

// Command represents a command in a state machine.
type Command[S, C any] = prop.Command[S, C]

// TestStateMachine tests a state machine using property-based testing.
// It generates sequences of commands and validates that the state machine
// behaves correctly according to the defined commands and their pre/post conditions.
//
// Parameters:
//   - t: The testing.T instance for the current test
//   - sm: The state machine to test
//   - cfg: Configuration for the property-based test
//
// Example:
//
//	sm := propx.StateMachine[BankAccount, BankCommand]{
//		InitialState: BankAccount{Balance: 0},
//		Commands: []propx.Command[BankAccount, BankCommand]{
//			depositCmd,
//			withdrawCmd,
//			closeCmd,
//		},
//	}
//	propx.TestStateMachine(t, sm, propx.Default())
func TestStateMachine[S, C any](t *testing.T, sm StateMachine[S, C], cfg Config) {
	prop.TestStateMachine(t, sm, cfg)
}

// =============================================================================
// GENERATORS
// =============================================================================

// Generator is the interface that all generators must implement.
type Generator[T any] = gen.Generator[T]

// Size controls the scale and limits of generators.
type Size = gen.Size

// Shrinker proposes "smaller" candidates during the shrinking process.
type Shrinker[T any] = gen.Shrinker[T]

// Shrinking strategy constants
const (
	ShrinkStrategyBFS = gen.ShrinkStrategyBFS // breadth-first search
	ShrinkStrategyDFS = gen.ShrinkStrategyDFS // depth-first search
)

// SetShrinkStrategy sets the shrinking strategy for all generators.
func SetShrinkStrategy(s string) {
	gen.SetShrinkStrategy(s)
}

// GetShrinkStrategy returns the current shrinking strategy.
func GetShrinkStrategy() string {
	return gen.GetShrinkStrategy()
}

// =============================================================================
// BASIC GENERATORS
// =============================================================================

// Int generates random integers with automatic range based on Size.
func Int(size gen.Size) gen.Generator[int] {
	return gen.Int(size)
}

// IntRange generates random integers within a specified range.
func IntRange(min, max int) gen.Generator[int] {
	return gen.IntRange(min, max)
}

// Uint generates random unsigned integers with automatic range based on Size.
func Uint(size gen.Size) gen.Generator[uint] {
	return gen.Uint(size)
}

// UintRange generates random unsigned integers within a specified range.
func UintRange(min, max uint) gen.Generator[uint] {
	return gen.UintRange(min, max)
}

// String generates random strings using an alphabet and Size.
func String(alphabet string, size gen.Size) gen.Generator[string] {
	return gen.String(alphabet, size)
}

// StringAlpha generates strings using only alphabetic characters.
func StringAlpha(size gen.Size) gen.Generator[string] {
	return gen.StringAlpha(size)
}

// StringAlphaNum generates strings using alphanumeric characters.
func StringAlphaNum(size gen.Size) gen.Generator[string] {
	return gen.StringAlphaNum(size)
}

// StringDigits generates strings using only digits.
func StringDigits(size gen.Size) gen.Generator[string] {
	return gen.StringDigits(size)
}

// StringASCII generates strings using all printable ASCII characters.
func StringASCII(size gen.Size) gen.Generator[string] {
	return gen.StringASCII(size)
}

// Bool generates random boolean values.
func Bool() gen.Generator[bool] {
	return gen.Bool()
}

// =============================================================================
// SLICE GENERATORS
// =============================================================================

// SliceOf generates random slices of the given type.
func SliceOf[T any](g gen.Generator[T], size gen.Size) gen.Generator[[]T] {
	return gen.SliceOf(g, size)
}

// =============================================================================
// COMBINATOR GENERATORS
// =============================================================================

// OneOf randomly selects one of the provided generators.
func OneOf[T any](generators ...gen.Generator[T]) gen.Generator[T] {
	return gen.OneOf(generators...)
}

// Const always returns the same value (without shrinking).
func Const[T any](v T) gen.Generator[T] {
	return gen.Const(v)
}

// Map applies f: A -> B preserving shrinking (maps A's candidates).
func Map[A, B any](ga gen.Generator[A], f func(A) B) gen.Generator[B] {
	return gen.Map(ga, f)
}

// Filter keeps only values that satisfy pred.
func Filter[T any](g gen.Generator[T], pred func(T) bool, maxTries int) gen.Generator[T] {
	return gen.Filter(g, pred, maxTries)
}

// Bind (flatMap): the output generator depends on the value generated in A.
func Bind[A, B any](ga gen.Generator[A], f func(A) gen.Generator[B]) gen.Generator[B] {
	return gen.Bind(ga, f)
}

// =============================================================================
// CUSTOM GENERATORS
// =============================================================================

// From creates a Generator from a function that implements the Generator interface.
func From[T any](fn func(*rand.Rand, Size) (T, Shrinker[T])) gen.Generator[T] {
	return gen.From(fn)
}

// =============================================================================
// DOMAIN-SPECIFIC GENERATORS
// =============================================================================

// CPF generates valid Brazilian CPF (Cadastro de Pessoas Físicas) numbers.
// If masked is true, returns formatted CPF (e.g., "123.456.789-01").
// If masked is false, returns raw CPF (e.g., "12345678901").
func CPF(masked bool) gen.Generator[string] {
	return domain.CPF(masked)
}

// CPFAny generates CPF with random masking (50/50 chance).
func CPFAny() gen.Generator[string] {
	return domain.CPFAny()
}

// ValidCPF validates if a string is a valid CPF.
func ValidCPF(s string) bool {
	return domain.ValidCPF(s)
}

// MaskCPF formats a raw CPF with dots and dashes.
func MaskCPF(raw string) string {
	return domain.MaskCPF(raw)
}

// UnmaskCPF removes formatting from a CPF string.
func UnmaskCPF(s string) string {
	return domain.UnmaskCPF(s)
}

// =============================================================================
// TESTING UTILITIES
// =============================================================================

// Equal compares two values of the same type and fails the test if they are not equal.
// It uses go-cmp for deep comparison and provides detailed diff output when values differ.
func Equal[T any](t *testing.T, got, want T) {
	quick.Equal(t, got, want)
}