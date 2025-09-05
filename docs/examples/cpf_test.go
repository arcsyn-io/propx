//go:build examples
// +build examples

// Package examples demonstrates how to use the propx property-based testing library.
// These examples show various testing patterns and how the shrinking mechanism
// helps find minimal counterexamples when properties fail.
package examples

import (
	"testing"

	"arcsyn.io/propx"
)

// Test_CPF_AlwaysValid demonstrates a property-based test for CPF generation.
// This test verifies that all generated CPF numbers are valid according to
// the CPF validation algorithm, and that the UnmaskCPF function is idempotent.
func Test_CPF_AlwaysValid(t *testing.T) {
	cfg := propx.Default()
	propx.ForAll(t, cfg, propx.CPF(false))(func(t *testing.T, cpf string) {
		if !propx.ValidCPF(cpf) {
			t.Fatalf("valid CPF generated was rejected: %q", cpf)
		}
		n1 := propx.UnmaskCPF(cpf)
		n2 := propx.UnmaskCPF(n1)
		propx.Equal(t, n1, n2)
	})
}

// Test_CPF_MaskUnmaskRoundTrip demonstrates testing the round-trip property
// of CPF masking and unmasking operations. This test verifies that
// unmasking a masked CPF and then masking it again produces the same result.
func Test_CPF_MaskUnmaskRoundTrip(t *testing.T) {
	propx.ForAll(t, propx.Default(), propx.CPF(true))(func(t *testing.T, masked string) {
		raw := propx.UnmaskCPF(masked)
		back := propx.UnmaskCPF(propx.MaskCPF(raw))
		propx.Equal(t, raw, back)
	})
}

// Test_CPF_Any_Valid demonstrates testing CPFAny() generator which produces
// CPF numbers with random masking (50/50 chance of masked or unmasked).
// This test verifies that all generated CPF numbers are valid regardless of format.
func Test_CPF_Any_Valid(t *testing.T) {
	propx.ForAll(t, propx.Default(), propx.CPFAny())(func(t *testing.T, s string) {
		if !propx.ValidCPF(s) {
			t.Fatalf("valid CPF generated was rejected: %q", s)
		}
	})
}

// Test_CPF_Invalid demonstrates a property-based test that is designed to fail.
// This test expects all CPF numbers to start with '9', which is not true for
// valid CPF generation. This example shows how the shrinking mechanism will
// find a minimal counterexample when the property fails.
func Test_CPF_Invalid(t *testing.T) {
	cfg := propx.Default()
	propx.ForAll(t, cfg, propx.CPF(false))(func(t *testing.T, cpf string) {
		if cpf[0] != '9' {
			t.Fatalf("expected to start with 9, but got %q", cpf)
		}
	})
}
