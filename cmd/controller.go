package main

import (
	"fmt"
	"os"

	"github.com/quic-s/quics-client/pkg/app"

	"github.com/quic-s/quics-client/pkg/viper"
	"github.com/spf13/cobra"
)

const (
	SIpCommand      = "server-host"
	SIpShortCommand = "d"

	SPortCommand      = "server-port"
	SPortShortCommand = "p"

	MyPortCommand      = "rest-port"
	MyPortShortCommand = "r"
)

var (
	SIp    string
	SPort  string
	MyPort string
)

func init() {
	startCmd := StartCmd()

	startCmd.Flags().StringVarP(&MyPort, MyPortCommand, "r", "", "my Port for make connection")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(ShutdownCmd())
	rootCmd.AddCommand(RebootCmd())

}

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Quics Client Server ",
		Run: func(cmd *cobra.Command, args []string) {
			if MyPort != "" {
				viper.WriteViperEnvVariables("REST_SERVER_PORT", MyPort)
			}
			app.RestServerStart()

		},
	}
}

func RebootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reboot",
		Short: "reboot the server",
		Run: func(cmd *cobra.Command, args []string) {
			app.Reboot()
		},
	}
}

func ShutdownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "shutdown",
		Short: "shutdown the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bye")
			os.Exit(0)
		},
	}

}
