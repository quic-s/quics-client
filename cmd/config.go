package main

import (
	"log"

	"github.com/quic-s/quics-client/pkg/sync"
	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	HostCommand      = "host"
	HostShortCommand = "H"

	PortCommand      = "port"
	PortShortCommand = "P"
)

func init() {
	configCmd := ConfigCmd()

	ChangeServerConfig := ChangeServerConfig()
	ChangeServerConfig.Flags().StringVarP(&SIp, HostCommand, HostShortCommand, "", "server domain/Ip for make connection")
	ChangeServerConfig.Flags().StringVarP(&SPort, PortCommand, PortShortCommand, "", "server Port for make connection")

	configCmd.AddCommand(ChangeServerConfig)

	configCmd.AddCommand(ReadConfig())

	rootCmd.AddCommand(configCmd)

}

func ConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "only read and rewrite config of [quics client], **can not delete config**",
	}
}

func ChangeServerConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "change server config",
		Run: func(cmd *cobra.Command, args []string) {
			result := sync.ConfigServer(SIp, SPort)
			log.Println("quics-client : [Config] : ", result)
		},
	}
}

func ReadConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show configs of quics client",
		Run: func(cmd *cobra.Command, args []string) {
			raw := utils.ReadEnvFile()
			result := "quics-client : [Show Config] : \n"

			for _, item := range raw {
				result += item + "/n"
			}
			log.Println(result)
		},
	}
}
