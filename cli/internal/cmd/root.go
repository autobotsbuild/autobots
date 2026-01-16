package cmd

import (
	"github.com/SheaHawkins/AutoBots/internal/cmd/hello"
	"github.com/SheaHawkins/AutoBots/internal/cmd/shared"
	"github.com/spf13/cobra"
)

var Version = "dev"

func NewRoot(deps shared.Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "autobots",
		Version:      Version,
		Short:        "AutoBots - agents that know the rules of the road.",
		Long:         `AutoBots is an AI-powered platform automating the long and boring tasks of software.`,
		SilenceUsage: true,
	}

	// Add subcommands
	cmd.AddCommand(hello.NewHelloCmd(deps))

	return cmd
}
