package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	SubCmd = &cobra.Command{
		Use:   "shutdown",
		Short: "shutdown the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("shutdown the server")
			os.Exit(1)
		},
	}
)
