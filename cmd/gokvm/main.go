package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gokvm",
		Short: "Software KVM",
	}

	rootCmd.AddCommand(newClientCmd())
	rootCmd.AddCommand(newServerCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}

	os.Exit(0)
}
