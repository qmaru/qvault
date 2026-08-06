package cmd

import (
	"qvault/cmd/api"
	"qvault/cmd/dbs"
	"qvault/cmd/generate"
	"qvault/cmd/manage"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "qvault",
		Short:   "qmaru key management service",
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableFlagsInUseLine = true
	rootCmd.AddCommand(
		api.NewCmd(),
		dbs.NewCmd(),
		manage.NewCmd(),
		generate.NewCmd(),
	)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
