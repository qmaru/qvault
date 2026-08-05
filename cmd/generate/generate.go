package generate

import (
	"fmt"
	"log"

	"qkms/services/kms"
	"qkms/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "generate command",
	}

	cmd.AddCommand(masterKeyGenerate(), apiKeyGenerate(), userCreate())

	return cmd
}

func userCreate() *cobra.Command {
	var username string

	cmd := &cobra.Command{
		Use:   "user",
		Short: "create user",
		Run: func(cmd *cobra.Command, args []string) {
			key, err := kms.CreateUser(username)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(key)
		},
	}

	cmd.Flags().StringVarP(&username, "name", "n", "", "username")
	cmd.MarkFlagRequired("name")

	return cmd
}

func masterKeyGenerate() *cobra.Command {
	return &cobra.Command{
		Use:   "master",
		Short: "generate master key",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := utils.GenerateMasterKey()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(result)
		},
	}
}

func apiKeyGenerate() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "generate api key",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := utils.GenerateAPIKey()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("sk-" + result)
		},
	}
}
