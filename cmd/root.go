package cmd

import (
	"qvault/cmd/api"
	"qvault/cmd/dbs"
	"qvault/cmd/generate"
	"qvault/cmd/manage"
	"qvault/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "qvault",
		Short:   "qvault key management service",
		Version: utils.Version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = true
	cmd.DisableFlagsInUseLine = true
	cmd.AddCommand(
		api.NewCmd(),
		dbs.NewCmd(),
		manage.NewCmd(),
		generate.NewCmd(),
	)

	return cmd
}

func Execute() error {
	if err := NewCmd().Execute(); err != nil {
		return err
	}
	return nil
}
