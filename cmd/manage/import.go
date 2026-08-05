package manage

import (
	"log"

	"qkms/services/kms"

	"github.com/spf13/cobra"
)

func importCmd() *cobra.Command {
	var apiKey string
	var input string
	cmd := &cobra.Command{
		Use:   "import",
		Short: "import command",
		Run: func(cmd *cobra.Command, args []string) {
			err := kms.ImportFromDotenv(apiKey, input)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	cmd.Flags().StringVarP(&input, "input", "i", "kms.env", "Input file path")
	cmd.MarkFlagRequired("key")

	return cmd
}
