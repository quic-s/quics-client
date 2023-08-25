package main

import (
	"fmt"
	"os"
	"os/exec"

	qhttp "github.com/quic-s/quics-client/pkg/http"
	"github.com/quic-s/quics-client/pkg/utils"
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

	// startCmd.Flags().StringVarP(&SIp, SIpCommand, "d", "", "server domain/Ip for make connection")
	// startCmd.Flags().StringVarP(&SPort, SPortCommand, "p", "", "server Port for make connection")
	startCmd.Flags().StringVarP(&MyPort, MyPortCommand, "r", "", "my Port for make connection")

	// if err := startCmd.MarkFlagRequired(SIpCommand); err != nil {
	// 	fmt.Println(err)
	// }
	// if err := startCmd.MarkFlagRequired(SPortCommand); err != nil {
	// 	fmt.Println(err)
	// }

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(ShutdownCmd())
	rootCmd.AddCommand(RebootCmd())

}

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start Quics Client Server ",
		Run: func(cmd *cobra.Command, args []string) {
			// 현재 IP 주소 중 첫 번째 것을 선택한다
			// hostname, err := os.Hostname()
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// addrs, err := net.LookupIP(hostname)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			// ip := addrs[0].String()

			// 아이피와 포트를 출력한다
			if MyPort != "" {
				utils.WriteViperEnvVariables("PORT", MyPort)
			}
			// utils.WriteViperEnvVariables("SERVER_IP", SIp)
			// utils.WriteViperEnvVariables("SERVER_PORT", SPort)

			qhttp.RestServerStart()
		},
	}
}

func RebootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reboot",
		Short: "reboot the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\n\trebooting ...\n")
			// 현재 실행 중인 프로세스의 경로와 인자를 얻습니다.
			str, err := os.Executable()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//fmt.Println("prevPid : ", os.Getpid())
			newProcess := exec.Command(str)
			newProcess.Stdout = os.Stdout
			newProcess.Stderr = os.Stderr
			newProcess.Stdin = os.Stdin
			newProcess.Env = os.Environ()

			err = newProcess.Start()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			//fmt.Println("newPid : ", newProcess.Process.Pid)
			os.Exit(0)
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
