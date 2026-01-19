package contract

// v0.0.1 Contract schema
//
// apiVersion: autobots/v1alpha1
// kind: Contract
// metadata:
//   is_draft: true
//   labels: { ... }
// spec:
//   consumer: { component: ... }
//   provider: { component: ... }
//   surface:
//     kind: http
//     http: { method, path, auth: { scheme, scopes } }
//   assertions: [{id, text}, ...]
//   bindings:
//     tests: [{id, kind, path, required, covers}, ...]

type Contract struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   ContractMeta   `yaml:"metadata"`
	Spec       ContractSpec   `yaml:"spec"`
}

type ContractMeta struct {
	IsDraft bool              `yaml:"is_draft"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

type ContractSpec struct {
	Consumer   ContractParty    `yaml:"consumer"`
	Provider   ContractParty    `yaml:"provider"`
	Surface    ContractSurface  `yaml:"surface"`
	Assertions []Assertion      `yaml:"assertions"`
	Bindings   ContractBindings `yaml:"bindings"`
}

type ContractParty struct {
	Component string `yaml:"component"`
}

type ContractSurface struct {
	Kind string       `yaml:"kind"` // currently: http
	HTTP *HTTPSurface `yaml:"http,omitempty"`
}

type HTTPSurface struct {
	Method string    `yaml:"method"` // e.g. POST
	Path   string    `yaml:"path"`   // e.g. /v1/invoices/{invoiceId}/mark-paid
	Auth   *HTTPAuth `yaml:"auth,omitempty"`
}

type HTTPAuth struct {
	Scheme string   `yaml:"scheme,omitempty"` // e.g. bearer
	Scopes []string `yaml:"scopes,omitempty"`
}

type Assertion struct {
	ID   string `yaml:"id"`
	Text string `yaml:"text"`
}

type ContractBindings struct {
	Tests []TestBinding `yaml:"tests"`
}

type TestBinding struct {
	ID       string   `yaml:"id"`
	Kind     string   `yaml:"kind"` // postman | sql (v0.0.1)
	Path     string   `yaml:"path"`
	Required bool     `yaml:"required"`
	Covers   []string `yaml:"covers,omitempty"` // assertion IDs
}
