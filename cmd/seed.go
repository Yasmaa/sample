package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)


var seedCmd *cobra.Command



func init() {

	seedCmd = &cobra.Command{
		Use:   "seed",
		Short: "seed the database with records",
		Long:  `command to seed the database with records`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("seed")
		},
	}



	dbCmd.AddCommand(seedCmd)
}








