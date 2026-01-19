package shared

import "time"

// Contract represents the root schema of the contract definition.
// apiVersion: autobots/v1alpha1
// kind: Contract
type Contract struct {
	ApiVersion   string           `json:"apiVersion" yaml:"apiVersion"`
	Kind         string           `json:"kind" yaml:"kind"`
	Metadata     ContractMetadata `json:"metadata" yaml:"metadata"`
	Spec         ContractSpec     `json:"spec" yaml:"spec"`
	SpecRevision Checksum         `json:"specRevision" yaml:"specRevision"`
	Changelog    []ChangelogEntry `json:"changelog" yaml:"changelog"`
}

type ContractStatus string

const (
	ContractStatusDraft      ContractStatus = "DRAFT"
	ContractStatusLocked     ContractStatus = "LOCKED"
	ContractStatusSatisfied  ContractStatus = "SATISFIED"
	ContractStatusStale      ContractStatus = "STALE"
	ContractStatusDeprecated ContractStatus = "DEPRECATED"
)

type ContractMetadata struct {
	ID        string            `json:"id" yaml:"id"`
	Slug      string            `json:"slug" yaml:"slug"`
	Status    ContractStatus    `json:"status" yaml:"status"`
	CreatedAt time.Time         `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt" yaml:"updatedAt"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
}

type ContractSpec struct {
	Consumer   ServiceInfo `json:"consumer" yaml:"consumer"`
	Provider   ServiceInfo `json:"provider" yaml:"provider"`
	Surface    Surface     `json:"surface" yaml:"surface"`
	Assertions []Assertion `json:"assertions" yaml:"assertions"`
	Tests      []Test      `json:"tests" yaml:"tests"`
}

type ServiceInfo struct {
	Component string `json:"component" yaml:"component"`
	Role      string `json:"role" yaml:"role"`
}

type Surface struct {
	Kind string       `json:"kind" yaml:"kind"` // http | ui | database | event
	HTTP *HTTPSurface `json:"http,omitempty" yaml:"http,omitempty"`
}

type HTTPSurface struct {
	OpenAPIRef string `json:"openapiRef" yaml:"openapiRef"`
}

type AssertionPriority string

const (
	AssertionPriorityMust   AssertionPriority = "MUST"
	AssertionPriorityShould AssertionPriority = "SHOULD"
	AssertionPriorityCould  AssertionPriority = "COULD"
)

type AssertionType string

const (
	AssertionTypeBehavior      AssertionType = "BEHAVIOR"
	AssertionTypeSchema        AssertionType = "SCHEMA"
	AssertionTypeInvariant     AssertionType = "INVARIANT"
	AssertionTypeError         AssertionType = "ERROR"
	AssertionTypeNonFunctional AssertionType = "NONFUNCTIONAL"
)

type Assertion struct {
	ID       string            `json:"id" yaml:"id"`
	Priority AssertionPriority `json:"priority" yaml:"priority"`
	Type     AssertionType     `json:"type" yaml:"type"`
	Text     string            `json:"text" yaml:"text"`
	Tags     []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Test struct {
	ID             string   `json:"id" yaml:"id"`
	Kind           string   `json:"kind" yaml:"kind"` // postman | sql
	Path           string   `json:"path" yaml:"path"`
	Covers         []string `json:"covers" yaml:"covers"`
	Required       bool     `json:"required" yaml:"required"`
	TimeoutSeconds int      `json:"timeoutSeconds" yaml:"timeoutSeconds"`
	Env            TestEnv  `json:"env" yaml:"env"`
	Fingerprint    Checksum `json:"fingerprint" yaml:"fingerprint"`
	Checks         []Check  `json:"checks" yaml:"checks"`
}

type TestEnv struct {
	RequiredVars []string `json:"requiredVars" yaml:"requiredVars"`
	OptionalVars []string `json:"optionalVars" yaml:"optionalVars"`
}

type Checksum struct {
	Algorithm string `json:"algorithm" yaml:"algorithm"`
	Value     string `json:"value" yaml:"value"`
}

type Check struct {
	Kind           string   `json:"kind" yaml:"kind"`
	Description    string   `json:"description" yaml:"description"`
	SchemaPath     string   `json:"schemaPath,omitempty" yaml:"schemaPath,omitempty"`
	MigrationsPath string   `json:"migrationsPath,omitempty" yaml:"migrationsPath,omitempty"`
	Covers         []string `json:"covers" yaml:"covers"`
}

type ChangelogEntry struct {
	At                time.Time `json:"at" yaml:"at"`
	By                Actor     `json:"by" yaml:"by"`
	Summary           string    `json:"summary" yaml:"summary"`
	Reason            string    `json:"reason" yaml:"reason"`
	AffectsAssertions []string  `json:"affectsAssertions" yaml:"affectsAssertions"`
	Link              Link      `json:"link" yaml:"link"`
}

type Actor struct {
	Type string `json:"type" yaml:"type"` // agent | human
	Name string `json:"name" yaml:"name"`
}

type Link struct {
	Kind string `json:"kind" yaml:"kind"` // github_pr
	Ref  string `json:"ref" yaml:"ref"`
}
