package dbs

import (
	"log"

	"qkms/dbs"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dbs",
		Short: "database manager",
	}

	cmd.AddCommand(
		createDB(),
	)

	return cmd
}

func createDB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a database",
		Run: func(cmd *cobra.Command, args []string) {
			err := dbs.CreateDB()
			if err != nil {
				log.Fatal(err)
			}

			err = dbs.CreateIndexes()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
