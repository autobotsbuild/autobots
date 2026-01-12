package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "autobots",
	Short: "AutoBots - AI-powered contract-driven development",
	Long: `AutoBots is an AI-powered platform for achieving 
contract-driven development at scale.

This CLI reconciles desired state from configuration files
and orchestrates operations through backend services.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", 
		"config file (default is $HOME/.autobots.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, 
		"enable verbose output")
	
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
		viper.SetConfigName(".autobots")
	}
	
	viper.AutomaticEnv()
	viper.ReadInConfig() // Ignore error if config doesn't exist
}
