package api

import (
	"log"

	"qkms/apis"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "api command",
		Run: func(cmd *cobra.Command, args []string) {
			err := apis.Run()
			if err != nil {
				log.Fatalf("failed to run api: %v\n", err)
			}
		},
	}
}
