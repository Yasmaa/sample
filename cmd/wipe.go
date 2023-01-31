package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)


var wipeCmd *cobra.Command



func init() {

	wipeCmd = &cobra.Command{
		Use:   "wipe",
		Short: "command to wipe your database",
		Long:  `command to drop all tables, views and types of the database`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("wipe")
		},
	}



	dbCmd.AddCommand(wipeCmd)
}








