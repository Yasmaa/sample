package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)


var backupCmd *cobra.Command



func init() {

	backupCmd = &cobra.Command{
		Use:   "backup",
		Short: "backup the database",
		Long:  `Command to backup the database`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("backup")
		},
	}



	dbCmd.AddCommand(backupCmd)
}








