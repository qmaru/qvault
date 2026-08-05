package export

import (
	"log"

	"qkms/services/kms"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	var apiKey string
	var output string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export command",
		Run: func(cmd *cobra.Command, args []string) {
			err := kms.ExportToDotenv(apiKey, output)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "k", "", "API Key")
	cmd.Flags().StringVarP(&output, "output", "o", "kms.env", "Output file path")
	cmd.MarkFlagRequired("key")

	return cmd
}
