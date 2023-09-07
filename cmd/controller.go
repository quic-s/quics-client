package main

import (
	"fmt"
	"os"

	"github.com/quic-s/quics-client/pkg/app"
	"github.com/quic-s/quics-client/pkg/http3"

	"github.com/quic-s/quics-client/pkg/viper"
	"github.com/spf13/cobra"
)

const (
	MyPortCommand      = "hport"
	MyPortShortCommand = "p"

	StartConnectPortCommand = "port"
)

var (
	MyPort string

	StartConnectPort       string
	StartConnectServerIp   string
	StartConnectServerport string
)

func init() {
	startCmd := StartCmd()

	startCmd.Flags().StringVarP(&MyPort, MyPortCommand, "r", "", "my Port for make connection")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(ShutdownCmd())
	rootCmd.AddCommand(RebootCmd())

}
func StartConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "start Client server and connect with Quics Server ",
		Run: func(cmd *cobra.Command, args []string) {
			// if MyPort != "" {
			// 	viper.WriteViperEnvVariables("REST_SERVER_PORT", MyPort)
			// }
			// http3.RestServerStart()

		},
	}
}

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Quics Client Server ",
		Run: func(cmd *cobra.Command, args []string) {
			if MyPort != "" {
				viper.WriteViperEnvVariables("REST_SERVER_PORT", MyPort)
			}
			http3.RestServerStart()

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
