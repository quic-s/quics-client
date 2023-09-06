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

	DirNNCommand      = "name"
	DirNNShortCommand = "n"

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

	ChangeRootDirConfig := ChangeRootDirConfig()
	ChangeRootDirConfig.Flags().StringVarP(&DirAbsPath, DirAbsPathCommand, DirAbsPathShortCommand, "", "directory absolute path")
	ChangeRootDirConfig.Flags().StringVarP(&DirNN, DirNNCommand, DirNNShortCommand, "", "directory nickname")

	if err := ChangeRootDirConfig.MarkFlagRequired(DirAbsPathCommand); err != nil {
		log.Println(err)
	}
	if err := ChangeRootDirConfig.MarkFlagRequired(DirNNCommand); err != nil {
		log.Println(err)
	}

	configCmd.AddCommand(ChangeRootDirConfig)
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

func ChangeRootDirConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "root",
		Short: "change root directory config",
		Run: func(cmd *cobra.Command, args []string) {
			//TODO 중복된 키, 중복된 value 금지
			if utils.IsDuplicateKey(DirNN) {
				log.Fatal("duplicate key, cannot overwrite")
			}
			if utils.IsDuplicateValue(DirAbsPath) {
				log.Fatal("duplicate value, cannot overwrite")
			}

			viper.WriteViperEnvVariables("ROOT."+DirNN, DirAbsPath)
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
