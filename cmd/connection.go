package main

import (
	"encoding/json"
	"log"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/spf13/cobra"
)

const (
	PasswordCommand      = "password"
	PasswordShortCommand = "p"

	ConnectLocalCommand      = "local"
	ConnectLocalShortCommand = "l"

	ConnectRemoteCommand      = "remote"
	ConnectRemoteShortCommand = "r"

	DisconnectRootCommand      = "root"
	DisconnectRootShortCommand = "r"
)

var (
	SIp                 string
	SPort               string
	UUID                string
	ClientPW            string
	DisConnectClientPW  string
	DisConnectRootDirPw string

	LocalRootDir      string
	RemoteRootDir     string
	DisconnectRootDir string
	RootDirPW         string
)

func init() {
	//connect
	ConnectCmd := ConnectCmd()
	ConnectServerCmd := ConnectServerCmd()
	ConnectRootCmd := ConnectRootCmd()

	ConnectServerCmd.Flags().StringVarP(&SIp, HostCommand, HostShortCommand, "", "server domain/Ip for make connection")
	ConnectServerCmd.Flags().StringVarP(&SPort, PortCommand, PortShortCommand, "", "server Port for make connection")
	ConnectServerCmd.Flags().StringVarP(&ClientPW, PasswordCommand, PasswordShortCommand, "", "password for entering server")
	if err := ConnectServerCmd.MarkFlagRequired(HostCommand); err != nil {
		log.Println(err)
	}

	ConnectCmd.AddCommand(ConnectServerCmd)
	//ConnectCmd.AddCommand(PingCmd())

	ConnectRootCmd.Flags().StringVarP(&LocalRootDir, ConnectLocalCommand, ConnectLocalShortCommand, "", "decide local root directory")
	ConnectRootCmd.Flags().StringVarP(&RemoteRootDir, ConnectRemoteCommand, ConnectRemoteShortCommand, "", "decide remote root directory")
	ConnectRootCmd.Flags().StringVarP(&RootDirPW, PasswordCommand, PasswordShortCommand, "", "password for entering root dir")
	if err := ConnectRootCmd.MarkFlagRequired(PasswordCommand); err != nil {
		log.Println(err)
	}

	ConnectCmd.AddCommand(ConnectRootCmd)
	ConnectCmd.AddCommand(ShowRemoteRootListCmd())
	rootCmd.AddCommand(ConnectCmd)

	//disconnect
	DisconnectCmd := DisconnectCmd()
	DisconnectServerCmd := DisconnectServerCmd()
	DisconnectServerCmd.Flags().StringVarP(&DisConnectClientPW, PasswordCommand, PasswordShortCommand, "", "password for disconnect server")
	if err := DisconnectServerCmd.MarkFlagRequired(PasswordCommand); err != nil {
		log.Println(err)
	}
	DisconnectCmd.AddCommand(DisconnectServerCmd)

	DisconnectRootCmd := DisconnectRootCmd()
	DisconnectRootCmd.Flags().StringVarP(&DisconnectRootDir, DisconnectRootCommand, DisconnectRootShortCommand, "", "choose witch root be disable ")
	DisconnectRootCmd.Flags().StringVarP(&DisConnectRootDirPw, PasswordCommand, PasswordShortCommand, "", "password for disconnect root")
	if err := DisconnectRootCmd.MarkFlagRequired(DisconnectRootCommand); err != nil {
		log.Println(err)
	}
	if err := DisconnectRootCmd.MarkFlagRequired(PasswordCommand); err != nil {
		log.Println(err)
	}
	DisconnectCmd.AddCommand(DisconnectRootCmd)
	rootCmd.AddCommand(DisconnectCmd)

}

func ConnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "connect to server [options]",
	}
}

func ConnectServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "connect to server",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("host:", SIp)
			log.Println("port:", SPort)
			log.Println("password:", ClientPW)
			registerClient := &types.RegisterClientHTTP3{
				Host:     SIp,
				Port:     SPort,
				ClientPW: ClientPW,
			}

			body, err := json.Marshal(registerClient)
			if err != nil {
				log.Println(err)
			}
			restClient := NewRestClient()
			restClient.PostRequest("/api/v1/connect/server", "application/json", body)
		},
	}

}

func ConnectRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "root",
		Short: "make connection with root dir",
		Run: func(cmd *cobra.Command, args []string) {

			//TODO 경로외에 닉네임도 받을 수 있게 하기
			restClient := NewRestClient()

			if LocalRootDir != "" && RemoteRootDir == "" { // local to remote

				registerRootDirHTTP3 := &types.RegisterRootDirHTTP3{
					RootDir:   LocalRootDir,
					RootDirPw: RootDirPW,
				}

				body, err := json.Marshal(registerRootDirHTTP3)
				if err != nil {
					log.Println(err)
				}

				restClient.PostRequest("/api/v1/connect/root/local", "application/json", body)

			}
			if RemoteRootDir != "" && LocalRootDir == "" { // romote to local

				registerRootDirHTTP3 := &types.RegisterRootDirHTTP3{
					RootDir:   RemoteRootDir,
					RootDirPw: RootDirPW,
				}
				body, err := json.Marshal(registerRootDirHTTP3)
				if err != nil {
					log.Println(err)
				}
				restClient.PostRequest("/api/v1/connect/root/remote", "application/json", body)

			}

		},
	}
}

func ShowRemoteRootListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-remote",
		Short: "get remote root list",
		Run: func(cmd *cobra.Command, args []string) {

			restClient := NewRestClient()
			restClient.GetRequest("/api/v1/connect/list/remote")
		},
	}
}

func DisconnectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disconnect",
		Short: "disconnect [options]",
	}
}

func DisconnectServerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "disconnect server",
		Run: func(cmd *cobra.Command, args []string) {
			restClient := NewRestClient()
			disconnectClientHTTP3 := &types.RegisterClientHTTP3{
				ClientPW: DisConnectClientPW,
			}
			body, err := json.Marshal(disconnectClientHTTP3)
			if err != nil {
				log.Println(err)
			}
			restClient.PostRequest("/api/v1/disconnect/server", "application/json", body)

		},
	}
}

func DisconnectRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "root",
		Short: "disconnect root",
		Run: func(cmd *cobra.Command, args []string) {

			restClient := NewRestClient()
			disconnectRootHTTP3 := &types.RegisterRootDirHTTP3{
				RootDir:   DisconnectRootDir,
				RootDirPw: DisConnectRootDirPw,
			}
			body, err := json.Marshal(disconnectRootHTTP3)
			if err != nil {
				log.Println(err)
			}
			restClient.PostRequest("/api/v1/disconnect/root", "application/json", body)
		},
	}
}

// DISCUSS/ Is it needed?
// func PingCmd() *cobra.Command {
// 	return &cobra.Command{
// 		Use:  "ping",
// 		Long: "check connection with server",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			log.Println(connection.Ping())
// 		},
// 	}
// }
