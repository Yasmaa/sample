package cmd

import (
	"api/config"
	"github.com/spf13/cobra"

)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "Backup API",
		Short: "Databases Backup API",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

}

func initConfig() {
	config.LoadConfig()

}
