package hello

import "github.com/spf13/cobra"

type Flags struct {
	Verbose bool
}

func (f *Flags) Bind(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&f.Verbose, "verbose", "v", false, "enable verbose output")
}
