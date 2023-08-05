package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qic",
	Short: "qic is a CLI for interacting with the quics client",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {

}
