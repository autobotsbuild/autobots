package cmd

import (
	"testing"

	"github.com/autobotsbuild/autobots/cmd/shared"
	"github.com/stretchr/testify/assert"
)

func TestNewRoot_Version(t *testing.T) {
	// Setup
	deps := shared.Dependencies{}
	originalVersion := Version
	defer func() { Version = originalVersion }()

	// Test Case 1: Default Version
	Version = "dev"
	cmd := NewRoot(deps)
	assert.Equal(t, "dev", cmd.Version)

	// Test Case 2: Injected Version
	Version = "v1.0.0"
	cmd = NewRoot(deps)
	assert.Equal(t, "v1.0.0", cmd.Version)
}
