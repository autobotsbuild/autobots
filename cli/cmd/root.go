package cmd

import (
	"github.com/autobotsbuild/autobots/cmd/hello"
	"github.com/autobotsbuild/autobots/cmd/shared"
	"github.com/spf13/cobra"
)

var Version = "dev"

func NewRoot(deps shared.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "autobots",
		Version:      Version,
		Short:        "autobots - agents that know the rules of software.",
		Long:         `autobots is an AI-powered platform automating the long and boring tasks of software.`,
		SilenceUsage: true,
	}

	// Add subcommands
	cmd.AddCommand(hello.NewHelloCmd(deps))

	return cmd
}
