package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/exec"
)

var (
	Server string
	SPort  string
	Host   string
	HPort  string
)

func init() {
	startCmd := StartCmd()

	startCmd.Flags().StringVarP(&Server, "server", "s", "", "server-IP for make connection")
	startCmd.Flags().StringVarP(&SPort, "port", "p", "", "server-Port for make connection")

	if err := startCmd.MarkFlagRequired("server"); err != nil {
		fmt.Println(err)
	}
	if err := startCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(ShutdownCmd())
	rootCmd.AddCommand(RebootCmd())

}

func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [options]",
		Short: "start Quics Client Server ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\n Start Quics-Client \n")

			hostname, err := os.Hostname()
			if err != nil {
				fmt.Println(err)
				return
			}

			addrs, err := net.LookupIP(hostname)
			if err != nil {
				fmt.Println(err)
				return
			}
			// IP 주소 중 첫 번째 것을 선택한다
			ip := addrs[0].String()

			// port := os.Getenv("PORT")
			// if port == "" {
			// 	// 환경 변수 PORT가 없으면, 기본값으로 8080을 사용한다
			// 	port = "8080"
			// }

			// 현재 아이피와 포트를 출력한다
			fmt.Println("Local IP    ", ip)
			fmt.Println("Server IP   ", Server)
			fmt.Println("Server Port ", SPort)

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
