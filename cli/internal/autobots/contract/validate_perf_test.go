package contract

import (
	"testing"
	"time"
)

// BenchmarkValidateContract is the idiomatic Go way to measure performance.
// Run this with: go test -bench=. -benchmem
func BenchmarkValidateContract(b *testing.B) {
	// Setup a clean contract once
	c := getPerfContract()

	// Reset timer to ignore setup costs
	b.ResetTimer()

	// The loop that Go's testing tool controls
	for i := 0; i < b.N; i++ {
		// We ignore the error because we are profiling the *cost* of the check,
		// not correctness (which is covered in unit tests).
		_ = ValidateContract(c)
	}
}

// TestValidationLatency enforces a strict SLA on validation logic.
// This ensures no one accidentally introduces a heavy operation (like a Regex compilation inside a loop or a network call)
// into the validation hot path.
func TestValidationLatency(t *testing.T) {
	c := getPerfContract()

	// We run the validation N times to average out OS scheduler noise.
	// A single run is too fast to measure reliably (nanoseconds).
	iterations := 5000
	threshold := 500 * time.Microsecond // 0.5ms per validation (extremely generous buffer)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		err := ValidateContract(c)
		if err != nil {
			t.Fatalf("Performance test failed: input contract became invalid: %v", err)
		}
	}
	totalTime := time.Since(start)
	avgTime := totalTime / time.Duration(iterations)

	// Report the actual time for visibility in test logs
	t.Logf("Average validation time: %v", avgTime)

	if avgTime > threshold {
		t.Errorf("Validation performance regression! Average time %v exceeded threshold %v", avgTime, threshold)
	}
}

// Helper to generate a slightly complex contract to simulate real load
func getPerfContract() *Contract {
	return &Contract{
		APIVersion: "autobots/v1alpha1",
		Kind:       "Contract",
		Metadata: ContractMeta{
			IsDraft: false,
		},
		Spec: ContractSpec{
			Consumer: ContractParty{Component: "web-ui"},
			Provider: ContractParty{Component: "billing-api"},
			Surface: ContractSurface{
				Kind: "http",
				HTTP: &HTTPSurface{
					Method: "POST",
					Path:   "/v1/invoices/very/long/path/to/test/string/allocation",
					Auth:   &HTTPAuth{Scheme: "Bearer"},
				},
			},
			Assertions: []Assertion{
				{ID: "a1", Text: "Status is 200"},
				{ID: "a2", Text: "Body is JSON"},
				{ID: "a3", Text: "Latency < 500ms"},
			},
			Bindings: ContractBindings{
				Tests: []TestBinding{
					{ID: "t1", Kind: "postman", Path: "./tests/t1.json", Required: true, Covers: []string{"a1", "a2"}},
					{ID: "t2", Kind: "sql", Path: "./tests/t2.sql", Required: false, Covers: []string{"a3"}},
				},
			},
		},
	}
}