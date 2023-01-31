package cmd

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "db cmd for database management",
	Long:  `db cmd for database management`,
}

func init() {

	rootCmd.AddCommand(dbCmd)
}
