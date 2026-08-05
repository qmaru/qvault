package main

import (
	"log"

	"qkms/cmd"
	"qkms/utils"
)

func main() {
	err := utils.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}

	err = cmd.Execute()
	if err != nil {
		log.Fatalf("failed to execute command: %v", err)
	}
}
