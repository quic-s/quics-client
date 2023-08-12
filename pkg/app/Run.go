package app

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var (
	SubCmd = &cobra.Command{
		Use:   "shutdown",
		Short: "shutdown the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.println("shutdown the server")
			os.Exit(1)
		}	
	}
)

init(){
	
}