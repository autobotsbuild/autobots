package hello

import (
	"github.com/SheaHawkins/AutoBots/internal/cmd/shared"
	"github.com/spf13/cobra"
)

func NewHelloCmd(deps shared.Dependencies) *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Use:   "hello [name]",
		Short: "Say hello",
		Long:  `A simple hello command to verify the CLI is working correctly.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd.Context(), deps, flags, args)
		},
	}
	flags.Bind(cmd)
	return cmd
}
