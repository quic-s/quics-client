package main

import (
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(SyncCmd())

}

func SyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "make sync the server",
	}
}
