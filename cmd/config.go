package main

import (
	"log"

	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
	"github.com/spf13/cobra"
)

const (
	HostCommand      = "host"
	HostShortCommand = "H"

	PortCommand      = "port"
	PortShortCommand = "P"

	DirAbsPathCommand      = "abspath"
	DirAbsPathShortCommand = "d"

	DirNNDeleteCommand      = "key"
	DirNNDeleteShortCommand = "k"
)

var (
	DirAbsPath string
	DirNN      string
	Key        string
)

func init() {
	configCmd := ConfigCmd()

	ChangeServerConfig := ChangeServerConfig()
	ChangeServerConfig.Flags().StringVarP(&SIp, HostCommand, HostShortCommand, "", "server domain/Ip for make connection")
	ChangeServerConfig.Flags().StringVarP(&SPort, PortCommand, PortShortCommand, "", "server Port for make connection")

	if err := ChangeServerConfig.MarkFlagRequired(HostCommand); err != nil {
		log.Println(err)
	}

	configCmd.AddCommand(ChangeServerConfig)

	configCmd.AddCommand(ReadConfig())

	DeleteConfig := DeleteConfig()
	DeleteConfig.Flags().StringVarP(&Key, DirNNDeleteCommand, DirNNDeleteShortCommand, "", "delete config by key")

	if err := DeleteConfig.MarkFlagRequired(DirNNDeleteCommand); err != nil {
		log.Println(err)
	}
	configCmd.AddCommand(DeleteConfig)

	rootCmd.AddCommand(configCmd)
}

func ConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "read config quics client",
	}
}

func ChangeServerConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "change server config",
		Run: func(cmd *cobra.Command, args []string) {
			viper.WriteViperEnvVariables("QUICS_SERVER_HOST", SIp)
			if SPort == "" {
				viper.WriteViperEnvVariables("QUICS_SERVER_PORT", SPort)
			}
		},
	}
}

func ReadConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "show configs of quics client",
		Run: func(cmd *cobra.Command, args []string) {
			utils.ReadEnvFile()
		},
	}
}

func DeleteConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "delete config of quics client",
		Run: func(cmd *cobra.Command, args []string) {
			viper.DeleteViperVariablesByKey(Key)
		},
	}
}
