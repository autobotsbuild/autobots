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

// ValidateContract validates the v0.0.1 schema and internal references.
// It does NOT enforce test pass/fail (that’s CI/runner responsibility).
func ValidateContract(c *Contract) error {
	var errs ValidationErrors

	// Top-level required
	if strings.TrimSpace(c.APIVersion) == "" {
		errs = append(errs, ValidationError{Path: "apiVersion", Msg: "required"})
	}
	if strings.TrimSpace(c.Kind) == "" {
		errs = append(errs, ValidationError{Path: "kind", Msg: "required"})
	} else if c.Kind != "Contract" {
		errs = append(errs, ValidationError{Path: "kind", Msg: "must be 'Contract'"})
	}

	// Metadata required: is_draft must exist in YAML, but in Go it’s bool default false.
	// If you want to detect “missing vs false”, you’d use *bool. v0.0.1 doesn’t require that,
	// so we accept false as a valid explicit value.
	// Labels optional.

	// Parties required
	if strings.TrimSpace(c.Spec.Consumer.Component) == "" {
		errs = append(errs, ValidationError{Path: "spec.consumer.component", Msg: "required"})
	}
	if strings.TrimSpace(c.Spec.Provider.Component) == "" {
		errs = append(errs, ValidationError{Path: "spec.provider.component", Msg: "required"})
	}

	// Surface required
	if strings.TrimSpace(c.Spec.Surface.Kind) == "" {
		errs = append(errs, ValidationError{Path: "spec.surface.kind", Msg: "required"})
	} else {
		switch c.Spec.Surface.Kind {
		case "http":
			if c.Spec.Surface.HTTP == nil {
				errs = append(errs, ValidationError{Path: "spec.surface.http", Msg: "required for kind=http"})
			} else {
				if strings.TrimSpace(c.Spec.Surface.HTTP.Method) == "" {
					errs = append(errs, ValidationError{Path: "spec.surface.http.method", Msg: "required"})
				}
				if strings.TrimSpace(c.Spec.Surface.HTTP.Path) == "" {
					errs = append(errs, ValidationError{Path: "spec.surface.http.path", Msg: "required"})
				} else if !strings.HasPrefix(c.Spec.Surface.HTTP.Path, "/") {
					errs = append(errs, ValidationError{Path: "spec.surface.http.path", Msg: "must start with '/'"})
				}
				if c.Spec.Surface.HTTP.Auth != nil {
					// scheme/scopes optional, but validate sane values if present
					if strings.ContainsAny(c.Spec.Surface.HTTP.Auth.Scheme, " \t\n") {
						errs = append(errs, ValidationError{Path: "spec.surface.http.auth.scheme", Msg: "must not contain whitespace"})
					}
				}
			}
		default:
			errs = append(errs, ValidationError{Path: "spec.surface.kind", Msg: "unsupported kind (v0.0.1 supports only 'http')"})
		}
	}

	// Assertions required
	if len(c.Spec.Assertions) == 0 {
		errs = append(errs, ValidationError{Path: "spec.assertions", Msg: "must have at least one assertion"})
	}
	assertionIDs := make(map[string]struct{}, len(c.Spec.Assertions))
	for i, a := range c.Spec.Assertions {
		prefix := fmt.Sprintf("spec.assertions[%d]", i)
		if strings.TrimSpace(a.ID) == "" {
			errs = append(errs, ValidationError{Path: prefix + ".id", Msg: "required"})
			continue
		}
		if _, ok := assertionIDs[a.ID]; ok {
			errs = append(errs, ValidationError{Path: prefix + ".id", Msg: "duplicate assertion id"})
		} else {
			assertionIDs[a.ID] = struct{}{}
		}
		if strings.TrimSpace(a.Text) == "" {
			errs = append(errs, ValidationError{Path: prefix + ".text", Msg: "required"})
		}
	}

	// Tests required (bindings.tests)
	if len(c.Spec.Bindings.Tests) == 0 {
		errs = append(errs, ValidationError{Path: "spec.bindings.tests", Msg: "must have at least one test binding"})
	}
	testIDs := make(map[string]struct{}, len(c.Spec.Bindings.Tests))
	for i, t := range c.Spec.Bindings.Tests {
		prefix := fmt.Sprintf("spec.bindings.tests[%d]", i)

		if strings.TrimSpace(t.ID) == "" {
			errs = append(errs, ValidationError{Path: prefix + ".id", Msg: "required"})
		} else {
			if _, ok := testIDs[t.ID]; ok {
				errs = append(errs, ValidationError{Path: prefix + ".id", Msg: "duplicate test id"})
			} else {
				testIDs[t.ID] = struct{}{}
			}
		}

		if strings.TrimSpace(t.Kind) == "" {
			errs = append(errs, ValidationError{Path: prefix + ".kind", Msg: "required"})
		} else {
			switch t.Kind {
			case "postman", "sql":
				// ok
			default:
				errs = append(errs, ValidationError{Path: prefix + ".kind", Msg: "unsupported kind (v0.0.1 supports postman|sql)"})
			}
		}

		if strings.TrimSpace(t.Path) == "" {
			errs = append(errs, ValidationError{Path: prefix + ".path", Msg: "required"})
		}

		// Covers must reference existing assertion IDs.
		for j, aid := range t.Covers {
			if _, ok := assertionIDs[aid]; !ok {
				errs = append(errs, ValidationError{
					Path: fmt.Sprintf("%s.covers[%d]", prefix, j),
					Msg:  fmt.Sprintf("unknown assertion id '%s'", aid),
				})
			}
		}

		// Optional: enforce that required tests cover something.
		// Comment out if you don’t want this rule.
		if t.Required && len(t.Covers) == 0 {
			errs = append(errs, ValidationError{Path: prefix + ".covers", Msg: "required tests should cover at least one assertion"})
		}
	}

	return errs.AsError()
}
