package main

import (
	"github.com/quic-s/quics-client/pkg/connection"
	"github.com/quic-s/quics-client/pkg/viper"
	"github.com/spf13/cobra"
)

func init() {
	ConnectCmd := ConnectCmd()
	ConnectCmd.Flags().StringVarP(&SIp, HostCommand, HostShortCommand, "", "server domain/Ip for make connection")
	ConnectCmd.Flags().StringVarP(&SPort, PortCommand, PortShortCommand, "", "server Port for make connection")
	rootCmd.AddCommand(ConnectCmd)

}

func ConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "connect to server",
		Run: func(cmd *cobra.Command, args []string) {
			if SIp == "" {
				SIp = viper.GetViperEnvVariables("QUICS_SERVER_IP")
			}
			if SPort == "" {
				SPort = viper.GetViperEnvVariables("QUICS_SERVER_PORT")
			}
			connection.RegisterQuicsClient()
		},
	}
}
