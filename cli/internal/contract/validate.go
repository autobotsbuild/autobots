package contract

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidContract = errors.New("invalid contract")
)

type ValidationError struct {
	Path string
	Msg  string
}

func (e ValidationError) Error() string {
	if e.Path == "" {
		return e.Msg
	}
	return fmt.Sprintf("%s: %s", e.Path, e.Msg)
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("contract validation failed:\n")
	for _, e := range ve {
		b.WriteString(" - ")
		b.WriteString(e.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func (ve ValidationErrors) AsError() error {
	if len(ve) == 0 {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrInvalidContract, ve.Error())
}

func (ve *ValidationErrors) Add(path, msg string) {
	*ve = append(*ve, ValidationError{Path: path, Msg: msg})
}

func (ve *ValidationErrors) Addf(path, format string, args ...any) {
	*ve = append(*ve, ValidationError{Path: path, Msg: fmt.Sprintf(format, args...)})
}

// ----------------------------
// Public entry point
// ----------------------------

func ValidateContract(c *Contract) error {
	var errs ValidationErrors

	if c == nil {
		errs.Add("", "contract is nil")
		return errs.AsError()
	}

	validateTopLevel(c, &errs)
	validateSpec(&c.Spec, &errs)

	return errs.AsError()
}

// ----------------------------
// Top-level validation
// ----------------------------

func validateTopLevel(c *Contract, errs *ValidationErrors) {
	if strings.TrimSpace(c.APIVersion) == "" {
		errs.Add("apiVersion", "required")
	}
	if strings.TrimSpace(c.Kind) == "" {
		errs.Add("kind", "required")
	} else if c.Kind != "Contract" {
		errs.Add("kind", "must be 'Contract'")
	}

	// metadata.is_draft is a bool; absence vs false is not distinguishable unless you use *bool.
	// v0.0.1 accepts false as a valid explicit value.
}

// ----------------------------
// Spec validation
// ----------------------------

func validateSpec(spec *ContractSpec, errs *ValidationErrors) {
	if spec == nil {
		errs.Add("spec", "required")
		return
	}

	validateParties(spec, errs)
	validateSurface(&spec.Surface, errs)

	assertionIDs := validateAssertions(spec.Assertions, errs)
	validateBindings(&spec.Bindings, assertionIDs, errs)
}

func validateParties(spec *ContractSpec, errs *ValidationErrors) {
	if strings.TrimSpace(spec.Consumer.Component) == "" {
		errs.Add("spec.consumer.component", "required")
	}
	if strings.TrimSpace(spec.Provider.Component) == "" {
		errs.Add("spec.provider.component", "required")
	}
}

// ----------------------------
// Surface validation
// ----------------------------

func validateSurface(s *ContractSurface, errs *ValidationErrors) {
	if s == nil {
		errs.Add("spec.surface", "required")
		return
	}

	if strings.TrimSpace(s.Kind) == "" {
		errs.Add("spec.surface.kind", "required")
		return
	}

	switch s.Kind {
	case "http":
		validateHTTP(s.HTTP, errs)
	default:
		errs.Add("spec.surface.kind", "unsupported kind (v0.0.1 supports only 'http')")
	}
}

func validateHTTP(h *HTTPSurface, errs *ValidationErrors) {
	if h == nil {
		errs.Add("spec.surface.http", "required for kind=http")
		return
	}

	if strings.TrimSpace(h.Method) == "" {
		errs.Add("spec.surface.http.method", "required")
	}
	if strings.TrimSpace(h.Path) == "" {
		errs.Add("spec.surface.http.path", "required")
	} else if !strings.HasPrefix(h.Path, "/") {
		errs.Add("spec.surface.http.path", "must start with '/'")
	}

	if h.Auth != nil {
		// auth fields optional; validate basic sanity if present
		if strings.ContainsAny(h.Auth.Scheme, " \t\n") {
			errs.Add("spec.surface.http.auth.scheme", "must not contain whitespace")
		}
		// scopes may be empty; no further constraints in v0.0.1
	}
}

// ----------------------------
// Assertions validation
// ----------------------------

type AssertionIDSet map[string]struct{}

func validateAssertions(assertions []Assertion, errs *ValidationErrors) AssertionIDSet {
	if len(assertions) == 0 {
		errs.Add("spec.assertions", "must have at least one assertion")
		return AssertionIDSet{}
	}

	ids := make(AssertionIDSet, len(assertions))
	for i, a := range assertions {
		prefix := fmt.Sprintf("spec.assertions[%d]", i)

		if strings.TrimSpace(a.ID) == "" {
			errs.Add(prefix+".id", "required")
			continue
		}
		if _, exists := ids[a.ID]; exists {
			errs.Add(prefix+".id", "duplicate assertion id")
		} else {
			ids[a.ID] = struct{}{}
		}

		if strings.TrimSpace(a.Text) == "" {
			errs.Add(prefix+".text", "required")
		}
	}

	return ids
}

// ----------------------------
// Bindings/tests validation
// ----------------------------

func validateBindings(b *ContractBindings, assertionIDs AssertionIDSet, errs *ValidationErrors) {
	if b == nil {
		errs.Add("spec.bindings", "required")
		return
	}

	if len(b.Tests) == 0 {
		errs.Add("spec.bindings.tests", "must have at least one test binding")
		return
	}

	testIDs := make(map[string]struct{}, len(b.Tests))
	for i, t := range b.Tests {
		prefix := fmt.Sprintf("spec.bindings.tests[%d]", i)

		if strings.TrimSpace(t.ID) == "" {
			errs.Add(prefix+".id", "required")
		} else {
			if _, ok := testIDs[t.ID]; ok {
				errs.Add(prefix+".id", "duplicate test id")
			} else {
				testIDs[t.ID] = struct{}{}
			}
		}

		if strings.TrimSpace(t.Kind) == "" {
			errs.Add(prefix+".kind", "required")
		} else {
			switch t.Kind {
			case "postman", "sql":
				// ok (v0.0.1)
			default:
				errs.Add(prefix+".kind", "unsupported kind (v0.0.1 supports postman|sql)")
			}
		}

		if strings.TrimSpace(t.Path) == "" {
			errs.Add(prefix+".path", "required")
		}

		// covers must reference existing assertions 
		for j, aid := range t.Covers {
			if _, ok := assertionIDs[aid]; !ok {
				errs.Addf(fmt.Sprintf("%s.covers[%d]", prefix, j), "unknown assertion id '%s'", aid)
			}
		}

		// required tests should cover at least one assertion.
		if t.Required && len(t.Covers) == 0 {
			errs.Add(prefix+".covers", "required tests should cover at least one assertion")
		}
	}
}
