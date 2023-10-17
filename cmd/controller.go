package main

import (
	"fmt"
	"log"
	"os"

	"github.com/quic-s/quics-client/pkg/net/http3"

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

// e.g. qic start
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

// e.g. qic reboot
func RebootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reboot",
		Short: "reboot the server",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			bres := restClient.GetRequest("/api/v1/reboot")
			log.Println(bres.String())

		},
	}
}

// e.g. qic shutdown
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
