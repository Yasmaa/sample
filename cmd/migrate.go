package cmd

import (
	"fmt"
	"log"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var migrateCmd *cobra.Command

func init() {

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "apply migrations to the database",
		Long:  `command to appy migrations to the database: migrate -- < up | down >`,
		Run: func(cmd *cobra.Command, args []string) {

			up, _ := cmd.Flags().GetBool("up")
			down, _ := cmd.Flags().GetBool("down")
			if up {
				log.Println("Running migrate up command")
			
				db, err := sql.Open("postgres", "postgresql:/myuser:password@localhost:5432/goSample?sslmode=disable")

				if err != nil {
					log.Fatal(err)
				}
				driver, err := postgres.WithInstance(db, &postgres.Config{})
				if err != nil {
					log.Fatal(err)
				}
				m, err := migrate.NewWithDatabaseInstance(
					"file:/internal/infrastructure/datastore/migrations",
					"postgres", driver)
				if err != nil {
					fmt.Println("here")
					log.Fatal(err)
				}

				if err := m.Up(); err != nil {
					log.Fatal(err)
				}
			} else if down {
				log.Println("Running migrate down command")
				db, err := sql.Open("postgres", "postgresql:/myuser:password@localhost:5432/goSample?sslmode=disable")

				if err != nil {
					log.Fatal(err)
				}
				driver, err := postgres.WithInstance(db, &postgres.Config{})
				if err != nil {
					log.Fatal(err)
				}
				m, err := migrate.NewWithDatabaseInstance(
					"file:/internal/infrastructure/datastore/migrations",
					"postgres", driver)
				if err != nil {
					log.Fatal(err)
				}

				if err := m.Down(); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println(`migrate cmd is used for database migration: add flag < up | down >`)
			}
		},
	}

	dbCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("up", "u", false, "migrate up")
	migrateCmd.Flags().BoolP("down", "d", false, "migrate down")

}
