package manage

import (
	"log"

	"qkms/services/manager"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manage",
		Short: "kms manager",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(exportCmd(), importCmd(), listKeyCmd())

	return cmd
}

func importCmd() *cobra.Command {
	var apiKey string
	var input string
	var prefix string

	cmd := &cobra.Command{
		Use:   "import",
		Short: "import a user's keys from a dotenv file",
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.ImportFromDotenv(apiKey, input, prefix)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "api key")
	cmd.Flags().StringVarP(&input, "input", "i", "kms.env", "input file path")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "key prefix")
	cmd.MarkFlagRequired("key")

	return cmd
}

func exportCmd() *cobra.Command {
	var apiKey string
	var output string
	var prefix string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "export a user's keys to a dotenv file",
		Run: func(cmd *cobra.Command, args []string) {
			err := manager.ExportToDotenv(apiKey, output, prefix)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "api key")
	cmd.Flags().StringVarP(&output, "output", "o", "kms.env", "output file path")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "key prefix")
	cmd.MarkFlagRequired("key")

	return cmd
}

func listKeyCmd() *cobra.Command {
	var apiKey string
	var prefix string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list a user's keys",
		Run: func(cmd *cobra.Command, args []string) {
			keys, err := manager.ListKeys(apiKey, prefix)
			if err != nil {
				log.Fatal(err)
			}
			for _, key := range keys {
				log.Println(key)
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "api key")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", "", "key prefix")
	cmd.MarkFlagRequired("key")

	return cmd
}
