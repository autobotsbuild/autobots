package contract

import (
	"testing"
	"time"
)

// Run this with: go test -bench=. -benchmem
func BenchmarkValidateContract(b *testing.B) {
	c := getPerfContract()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ValidateContract(c)
	}
}

// TestValidationLatency enforces a strict SLA on validation logic.
// This ensures no one accidentally introduces a heavy operation (like a Regex compilation inside a loop or a network call)
// into the validation hot path.
func TestValidationLatency(t *testing.T) {
	c := getPerfContract()

	iterations := 5000
	threshold := 500 * time.Microsecond

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