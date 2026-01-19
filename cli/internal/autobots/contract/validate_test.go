package contract

import (
	"strings"
	"testing"
)

func TestValidateContract(t *testing.T) {
	// Helper to create a valid "base" contract to minimize boilerplate in tests
	validContract := func() *Contract {
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
						Path:   "/v1/invoices",
						Auth:   &HTTPAuth{Scheme: "Bearer"},
					},
				},
				Assertions: []Assertion{
					{ID: "a1", Text: "Status is 200"},
				},
				Bindings: ContractBindings{
					Tests: []TestBinding{
						{ID: "t1", Kind: "postman", Path: "./tests/t1.json", Required: true, Covers: []string{"a1"}},
					},
				},
			},
		}
	}

	tests := []struct {
		name        string
		modifier    func(c *Contract) // Function to modify the valid contract for specific test cases
		expectError bool
		errContains []string // substrings expected in the error message
	}{
		{
			name:        "Valid Contract",
			modifier:    func(c *Contract) {},
			expectError: false,
		},
		{
			name: "Missing Top-Level Fields",
			modifier: func(c *Contract) {
				c.APIVersion = ""
				c.Kind = ""
			},
			expectError: true,
			errContains: []string{"apiVersion: required", "kind: required"},
		},
		{
			name: "Invalid Kind Value",
			modifier: func(c *Contract) {
				c.Kind = "NotContract"
			},
			expectError: true,
			errContains: []string{"kind: must be 'Contract'"},
		},
		{
			name: "Missing Parties",
			modifier: func(c *Contract) {
				c.Spec.Consumer.Component = ""
				c.Spec.Provider.Component = ""
			},
			expectError: true,
			errContains: []string{"spec.consumer.component: required", "spec.provider.component: required"},
		},
		{
			name: "Surface Kind Unsupported",
			modifier: func(c *Contract) {
				c.Spec.Surface.Kind = "grpc" // not supported in v0.0.1
			},
			expectError: true,
			errContains: []string{"spec.surface.kind: unsupported kind"},
		},
		{
			name: "HTTP Surface Missing Method and Path",
			modifier: func(c *Contract) {
				c.Spec.Surface.HTTP.Method = ""
				c.Spec.Surface.HTTP.Path = "no-slash"
			},
			expectError: true,
			errContains: []string{"spec.surface.http.method: required", "spec.surface.http.path: must start with '/'"},
		},
		{
			name: "HTTP Surface Nil",
			modifier: func(c *Contract) {
				c.Spec.Surface.HTTP = nil
			},
			expectError: true,
			errContains: []string{"spec.surface.http: required for kind=http"},
		},
		{
			name: "Auth Scheme Whitespace",
			modifier: func(c *Contract) {
				c.Spec.Surface.HTTP.Auth.Scheme = "Bearer Token"
			},
			expectError: true,
			errContains: []string{"spec.surface.http.auth.scheme: must not contain whitespace"},
		},
		{
			name: "Duplicate Assertion IDs",
			modifier: func(c *Contract) {
				c.Spec.Assertions = append(c.Spec.Assertions, Assertion{ID: "a1", Text: "duplicate"})
			},
			expectError: true,
			errContains: []string{"spec.assertions[1].id: duplicate assertion id"},
		},
		{
			name: "Duplicate Test IDs",
			modifier: func(c *Contract) {
				c.Spec.Bindings.Tests = append(c.Spec.Bindings.Tests, TestBinding{
					ID: "t1", Kind: "sql", Path: "x.sql",
				})
			},
			expectError: true,
			errContains: []string{"spec.bindings.tests[1].id: duplicate test id"},
		},
		{
			name: "Test Covers Non-Existent Assertion",
			modifier: func(c *Contract) {
				c.Spec.Bindings.Tests[0].Covers = []string{"ghost-assertion"}
			},
			expectError: true,
			errContains: []string{"spec.bindings.tests[0].covers[0]: unknown assertion id 'ghost-assertion'"},
		},
		{
			name: "Required Test Covers Nothing",
			modifier: func(c *Contract) {
				c.Spec.Bindings.Tests[0].Required = true
				c.Spec.Bindings.Tests[0].Covers = []string{}
			},
			expectError: true,
			errContains: []string{"spec.bindings.tests[0].covers: required tests should cover at least one assertion"},
		},
		{
			name: "Invalid Test Kind",
			modifier: func(c *Contract) {
				c.Spec.Bindings.Tests[0].Kind = "curl"
			},
			expectError: true,
			errContains: []string{"spec.bindings.tests[0].kind: unsupported kind"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validContract()
			tt.modifier(c)

			err := ValidateContract(c)

			if !tt.expectError && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				// Check if the error is of the correct wrapping type
				// (Though checking the string output is usually sufficient for validation logic)
				if !strings.Contains(err.Error(), "invalid contract") {
					t.Errorf("expected error to wrap 'invalid contract', got: %v", err)
				}

				for _, msg := range tt.errContains {
					if !strings.Contains(err.Error(), msg) {
						t.Errorf("expected error message to contain %q, got: \n%v", msg, err.Error())
					}
				}
			}
		})
	}
}