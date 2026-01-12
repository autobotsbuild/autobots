package cli

import (
	"fmt"
	
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var helloCmd = &cobra.Command{
	Use:   "hello [name]",
	Short: "Say hello",
	Long:  `A simple hello command to verify the CLI is working correctly.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := "World"
		if len(args) > 0 {
			name = args[0]
		}
		
		if viper.GetBool("verbose") {
			fmt.Println("Running in verbose mode...")
		}
		
		fmt.Printf("Hello, %s! Welcome to AutoBots.\n", name)
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
}
